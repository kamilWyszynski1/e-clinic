package handler

import (
	"e-clinic/src/backend/clinic"
	"e-clinic/src/backend/models"
	"errors"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

const selectFeeBySpecialist = `SELECT id, fee_per_30_min FROM specialist_fee WHERE specialist = ? AND speciality = ?`

func (h Handler) CreateAppointment(a *clinic.Appointment) (*models.Appointment, int, error) {
	log := h.log.WithField("method", "CreateAppointment")

	if err := a.Validate(h.now()); err != nil {
		log.WithError(err).Error("invalid appointment")
		return nil, http.StatusBadRequest, err
	}

	free, err := h.isSpecialistFree(a.SpecialistID, &clinic.TimeRange{
		From: a.ScheduledAt,
		To:   a.ScheduledAt.Add(a.Duration.Duration),
	})
	if free {
		goto Insert
	}

	//tr, _, err := h.GetSpecialistFreeTime(a.SpecialistID, &clinic.TimeRange{
	//	From: a.ScheduledAt,
	//	To:   a.ScheduledAt.Add(a.Duration.Duration),
	//})
	//if err != nil {
	//	return nil, http.StatusInternalServerError, nil
	//}
	//
	//for _, r := range tr.Ranges {
	//	if a.ScheduledAt.After(r.From) && a.ScheduledAt.Add(a.Duration.Duration).Before(r.To) {
	//		goto Insert
	//	}
	//}
	log.Warn("specialist not available")
	return nil, http.StatusBadRequest, errors.New("specialist is not available then")

Insert:
	type fee struct {
		ID          uuid.UUID
		FeePer30Min float64
	}
	var f fee
	err = h.db.SelectBySql(selectFeeBySpecialist, a.SpecialistID, a.Speciality).LoadOne(&f)
	if err != nil {
		log.WithError(err).Error("failed to query specialist fee")
		return nil, http.StatusBadRequest, err
	}
	ap := models.Appointment{
		ID:            uuid.NewV4(),
		State:         models.ApoitntmentstateenumAccepted,
		Patient:       a.PatientID,
		SpecialistFee: f.ID,
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

func (h Handler) CreatePatient(p *clinic.PatientRequest) (*models.Patient, int, error) {
	log := h.log.WithField("method", "CreatePatient")

	if err := p.Validate(); err != nil {
		log.WithError(err).Error("invalid patient request")
		return nil, http.StatusBadRequest, nil
	}

	patient := &models.Patient{
		ID:          uuid.NewV4(),
		Name:        p.Name,
		Surname:     p.Surname,
		Email:       p.Email,
		PhoneNumber: "500400300",
		Age:         20,
		Gender:      p.Gender,
	}
	if err := patient.Insert(h.db); err != nil {
		log.WithError(err).Error("failed to insert patient")
		return nil, http.StatusInternalServerError, nil
	}
	return patient, http.StatusCreated, nil
}
