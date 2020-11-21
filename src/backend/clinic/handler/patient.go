package handler

import (
	"e-clinic/src/backend/clinic"
	"e-clinic/src/backend/models"
	"errors"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

const selectFeeBySpecialist = `SELECT id FROM specialist_fee WHERE specialist = ? AND speciality = ?`

func (h Handler) CreateAppointment(a *clinic.Appointment) (*models.Appointment, int, error) {
	log := h.log.WithField("method", "CreateAppointment")

	if err := a.Invalidate(h.now()); err != nil {
		log.WithError(err).Error("invalid appointment")
		return nil, http.StatusBadRequest, err
	}

	tr, _, err := h.GetFreeTime(a.SpecialistID, &clinic.TimeRange{
		From: a.ScheduledAt,
		To:   a.ScheduledAt.Add(a.Duration.Duration),
	})
	if err != nil {
		return nil, http.StatusInternalServerError, nil
	}

	for _, r := range tr.Ranges {
		if a.ScheduledAt.After(r.From) && a.ScheduledAt.Add(a.Duration.Duration).Before(r.To) {
			goto Insert
		}
	}
	log.Warn("specialist not available")
	return nil, http.StatusBadRequest, errors.New("specialist is not available then")

Insert:

	var specialistFeeID uuid.UUID
	if err := h.db.SelectBySql(selectFeeBySpecialist, a.SpecialistID, a.Speciality).LoadOne(&specialistFeeID); err != nil {
		log.WithError(err).Error("failed to query specialist fee")
		return nil, http.StatusBadRequest, err
	}
	ap := models.Appointment{
		ID:            uuid.NewV4(),
		State:         models.ApoitntmentstateenumCreated,
		Patient:       a.PatientID,
		SpecialistFee: specialistFeeID,
		ScheduledTime: a.ScheduledAt,
		Duration:      int(a.Duration.Seconds()),
	}
	if err := ap.Insert(h.db); err != nil {
		log.WithError(err).Error("failed to insert appointment")
		return nil, http.StatusInternalServerError, nil
	}

	form := models.AppointmentForm{
		ID:          uuid.NewV4(),
		Appointment: ap.ID,
		Comment:     a.PatientComment,
		Symptoms:    a.PatientSymptoms,
	}
	if err := form.Insert(h.db); err != nil {
		log.WithError(err).Error("failed to insert AppointmentForm")
		return nil, http.StatusInternalServerError, nil
	}

	return &ap, http.StatusOK, nil
}
