package handler

import (
	"e-clinic/src/backend/clinic"
	"e-clinic/src/backend/models"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	uuid "github.com/satori/go.uuid"
)

func (h Handler) GetFreeTime(id uuid.UUID, tr *clinic.TimeRange) (*clinic.TimeRanges, int, error) {
	return &clinic.TimeRanges{Ranges: []clinic.TimeRange{
		{
			From: time.Now().Add(-time.Hour * 4),
			To:   time.Now().Add(-time.Hour * 3),
		},
		{
			From: time.Now().Add(-time.Hour),
			To:   time.Now().Add(time.Hour * 24),
		},
	}}, http.StatusOK, nil
}

const (
	doesSpecialistMatchSQL = `SELECT specialist = ? FROM appointment WHERE id = ?`
	finishAppointment      = `UPDATE appointment SET = 'OK' WHERE id = ?`
)

func (h Handler) MakePrescription(p *clinic.Prescription) (int, error) {
	log := h.log.WithField("method", "MakePrescription")
	if err := p.Invalidate(); err != nil {
		log.WithError(err).Error("invalid request")
		return http.StatusBadRequest, err
	}

	var match bool
	if err := h.db.SelectBySql(doesSpecialistMatchSQL, p.SpecialistID, p.AppointmentID).LoadOne(&match); err != nil {
		log.WithError(err).Error("failed to check if appointment match")
		return http.StatusInternalServerError, nil
	} else if !match {
		return http.StatusBadRequest, errors.New("this specialist cannot make prescription for another appointment")
	}

	res := models.AppointmentResult{
		ID:           uuid.NewV4(),
		Appointment:  p.AppointmentID,
		Comment:      p.Comment,
		Prescription: p.Prescription,
	}
	if err := res.Insert(h.db); err != nil {
		log.WithError(err).Error("failed to insert appointment result")
		return http.StatusInternalServerError, nil
	}

	if _, err := h.db.UpdateBySql(finishAppointment, p.AppointmentID).Exec(); err != nil {
		log.WithError(err).Error("failed to update appointment")
		return http.StatusInternalServerError, nil
	}
	return http.StatusOK, nil
}

func (h Handler) AcceptAppointment(aID uuid.UUID) (int, error) {
	log := h.log.WithFields(logrus.Fields{
		"method":        "AcceptAppointment",
		"appointmentID": aID,
	})
	if err := ChangeAppointmentStatus(h.db, aID, models.ApoitntmentstateenumAccepted); err != nil {
		log.WithError(err).Error("failed to change appointment status")
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (h Handler) RejectAppointment(aID uuid.UUID) (int, error) {
	log := h.log.WithFields(logrus.Fields{
		"method":        "RejectAppointment",
		"appointmentID": aID,
	})
	if err := ChangeAppointmentStatus(h.db, aID, models.ApoitntmentstateenumRejected); err != nil {
		log.WithError(err).Error("failed to change appointment status")
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
