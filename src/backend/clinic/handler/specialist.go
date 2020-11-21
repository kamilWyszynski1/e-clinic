package handler

import (
	"e-clinic/src/backend/clinic"
	"e-clinic/src/backend/models"
	"errors"
	"net/http"
	"time"

	"github.com/gocraft/dbr"

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
	checkStateAndSpecialist = `SELECT sf.specialist = ? FROM specialist_fee sf
JOIN appointment a on sf.id = a.specialist_fee
WHERE a.state = 'OK' AND a.id = ? and a.scheduled_time + make_interval(0,0,0,0,0,0, duration) > now()`
)

func (h Handler) MakePrescription(p *clinic.Prescription) (int, error) {
	log := h.log.WithField("method", "MakePrescription")
	if err := p.Invalidate(); err != nil {
		log.WithError(err).Error("invalid request")
		return http.StatusBadRequest, err
	}

	var match bool
	if err := h.db.SelectBySql(checkStateAndSpecialist, p.SpecialistID, p.AppointmentID).LoadOne(&match); errors.Is(err, dbr.ErrNotFound) {
		log.Debug("appointment does not match")
		return http.StatusBadRequest, nil
	} else if err != nil {
		log.WithError(err).Error("failed to check if appointment match")
		return http.StatusInternalServerError, nil
	} else if !match {
		return http.StatusBadRequest, errors.New("specialist does not match with appointment")
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

	if err := ChangeAppointmentStatus(h.db, p.AppointmentID, models.ApoitntmentstateenumFinished); err != nil {
		log.WithError(err).Error("failed to change appointment status")
		return http.StatusInternalServerError, err
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
