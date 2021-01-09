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
	log.Info("call")

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
	aInfos := make([]*clinic.AppointmentInfo, 0)
	for _, a := range appointments {
		i, _, err := h.GetAppointment(a.ID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		aInfos = append(aInfos, i)

	}
	return &clinic.AppointmentList{
		Appointments: aInfos,
		Len:          len(aInfos),
	}, http.StatusOK, nil
}

const (
	changeAppointmentStatus = `UPDATE appointment SET state = ? WHERE id = ? RETURNING true`
	getAppointmentStatus    = `SELECT state FROM appointment WHERE id =?`
)

func ChangeAppointmentStatus(db dbr.SessionRunner, apID uuid.UUID, status models.Apoitntmentstateenum) error {
	var current models.Apoitntmentstateenum
	if err := db.UpdateBySql(getAppointmentStatus, apID).Load(&current); err != nil {
		return fmt.Errorf("failed to check status, %w", err)
	}

	options, ok := transitions[current]
	if !ok {
		return errors.New("transition not expected")
	}
	for _, o := range options {
		if o == status {
			goto change
		}
	}
	return errors.New("invalid transition")

change:
	var done bool
	if err := db.UpdateBySql(changeAppointmentStatus, status, apID).Load(&done); err != nil {
		return fmt.Errorf("failed to update status, %w", err)
	} else if !done {
		return errors.New("appointment status not changed")
	}
	return nil
}

var transitions = map[models.Apoitntmentstateenum][]models.Apoitntmentstateenum{
	models.ApoitntmentstateenumCreated: {
		models.ApoitntmentstateenumAccepted, models.ApoitntmentstateenumRejected,
	},
	models.ApoitntmentstateenumAccepted: {
		models.ApoitntmentstateenumToBePaid,
	},
	models.ApoitntmentstateenumToBePaid: {
		models.ApoitntmentstateenumNotPaid, models.ApoitntmentstateenumOk,
	},
	models.ApoitntmentstateenumOk: {
		models.ApoitntmentstateenumFinished,
	},
}

const (
	appointmentFormByAppointmentID   = `select comment, symptoms from appointment_form where appointment = ?`
	appointmentResultByAppointmentID = `select comment from appointment_result where appointment = ?`
	drugDosingByAppointmentID        = `select * from appointment_result_prescription arp
join appointment_result ar on arp.appointment_result = ar.id
where ar.appointment = ?`
)

func (h Handler) GetAppointment(aID uuid.UUID) (*clinic.AppointmentInfo, int, error) {
	log := h.log.WithField("method", "GetAppointment")
	a, err := models.AppointmentByID(h.db, aID)
	if err != nil {
		log.WithError(err).Error("failed to query appointment")
		return nil, http.StatusInternalServerError, err
	}

	var form *models.AppointmentForm
	err = h.db.SelectBySql(appointmentFormByAppointmentID, aID).LoadOne(&form)
	if err != nil && !errors.Is(err, dbr.ErrNotFound) {
		log.WithError(err).Error("failed to query form")
		return nil, http.StatusInternalServerError, err
	}

	var specialistComment string
	err = h.db.SelectBySql(appointmentResultByAppointmentID, aID).LoadOne(&specialistComment)
	if err != nil && !errors.Is(err, dbr.ErrNotFound) {
		log.WithError(err).Error("failed to query form")
		return nil, http.StatusInternalServerError, err
	}

	var drugs []clinic.DrugDose
	err = h.db.SelectBySql(drugDosingByAppointmentID, aID).LoadOne(&drugs)
	if err != nil && !errors.Is(err, dbr.ErrNotFound) {
		log.WithError(err).Error("failed to query form")
		return nil, http.StatusInternalServerError, err
	}

	return &clinic.AppointmentInfo{
		Appointment: a,
		Form:        form,
		Prescription: &clinic.Prescription{
			Comment: specialistComment,
			Drugs:   drugs,
		},
	}, http.StatusOK, nil
}

const (
	getSpecialists    = `SELECT id, name, surname FROM specialist`
	getSpecialistFees = `SELECT * FROM specialist_fee WHERE specialist = ?`
)

func (h Handler) GetSpecialists() (*clinic.SpecialistList, int, error) {
	log := h.log.WithField("method", "GetSpecialists")

	var specialists []models.Specialist
	_, err := h.db.SelectBySql(getSpecialists).Load(&specialists)
	if err != nil && !errors.Is(err, dbr.ErrNotFound) {
		log.WithError(err).Error("failed to query specialists")
		return nil, http.StatusInternalServerError, err
	}

	list := make([]clinic.SpecialistWithFee, 0)

	for _, s := range specialists {
		var fees []models.SpecialistFee
		_, err := h.db.SelectBySql(getSpecialistFees, s.ID).Load(&fees)
		if err != nil && !errors.Is(err, dbr.ErrNotFound) {
			log.WithError(err).Error("failed to query specialists")
			return nil, http.StatusInternalServerError, err
		}
		for _, f := range fees {
			list = append(list, clinic.SpecialistWithFee{
				SpecialistID: s.ID,
				FeeId:        f.ID,
				Name:         s.Name,
				Surname:      s.Surname,
				Speciality:   f.Speciality,
				FeePer30Min:  f.FeePer30Min,
			})
		}
	}
	return &clinic.SpecialistList{Specialists: list}, http.StatusOK, nil
}
