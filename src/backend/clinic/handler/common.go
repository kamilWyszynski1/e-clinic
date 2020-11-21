package handler

import (
	"e-clinic/src/backend/clinic"
	"e-clinic/src/backend/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gocraft/dbr"
	uuid "github.com/satori/go.uuid"
)

const (
	patientAppointmentsSQL = `SELECT * FROM appointment
WHERE patient = ? AND scheduled_time BETWEEN ? AND ?`

	specialistAppointmentsSQL = `SELECT * FROM appointment
WHERE specialist = ? AND scheduled_time BETWEEN ? AND ?`
)

func (h Handler) GetAppointments(r *clinic.AppointmentsRequest) (*clinic.AppointmentList, int, error) {
	log := h.log.WithField("method", "GetAppointments")

	var appointments []models.Appointment
	switch r.UserType {
	case clinic.UserTypePatient:
		_, err := h.db.SelectBySql(patientAppointmentsSQL, r.ID, r.Range.From, r.Range.To).Load(&appointments)
		if err != nil {
			log.WithError(err).Error("failed to query appointments")
			return nil, http.StatusInternalServerError, err
		}
	case clinic.UserTypeSpecialist:
		_, err := h.db.SelectBySql(specialistAppointmentsSQL, r.ID, r.Range.From, r.Range.To).Load(&appointments)
		if err != nil {
			log.WithError(err).Error("failed to query appointments")
			return nil, http.StatusInternalServerError, err
		}
	default:
		log.Warnf("unknown userType: %s", r.UserType)
		return nil, http.StatusBadRequest, errors.New("invalid userType")
	}
	if len(appointments) == 0 {
		return nil, http.StatusNoContent, nil
	}
	return &clinic.AppointmentList{
		Appointments: appointments,
		Len:          len(appointments),
	}, http.StatusOK, nil
}

const changeAppointmentStatus = `UPDATE appointment SET state = ? WHERE id = ? RETURNING true`

func ChangeAppointmentStatus(db dbr.SessionRunner, apID uuid.UUID, status models.Apoitntmentstateenum) error {
	var done bool
	if err := db.UpdateBySql(changeAppointmentStatus, status, apID).Load(&done); err != nil {
		return fmt.Errorf("failed to update status, %w", err)
	} else if !done {
		return errors.New("appointment status not changed")
	}
	return nil
}
