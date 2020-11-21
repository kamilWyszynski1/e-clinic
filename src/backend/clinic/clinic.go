package clinic

import (
	"e-clinic/src/backend/models"
	"encoding/json"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Assistant interface {
	// GetAppointments returns appointments in given time range
	GetFreeTime(id uuid.UUID, timeRange *TimeRange) (*TimeRanges, int, error)
	MakePrescription(p *Prescription) (int, error)
	AcceptAppointment(aID uuid.UUID) (int, error)
	RejectAppointment(aID uuid.UUID) (int, error)

	// CreateAppointment checks if appointment is valid and schedules it
	CreateAppointment(a *Appointment) (*models.Appointment, int, error)
	GetAppointments(ar *AppointmentsRequest) (*AppointmentList, int, error)
}

type AppointmentList struct {
	Appointments []models.Appointment `json:"appointments"`
	Len          int                  `json:"len"`
}

type Prescription struct {
	AppointmentID uuid.UUID `json:"appointment_id"`
	SpecialistID  uuid.UUID `json:"specialist_id"`
	Comment       string    `json:"comment"`
	Prescription  string    `json:"prescription"` // listed drugs
}

func (p Prescription) Invalidate() error {
	if p.Comment == "" {
		return errors.New("comment cannot be empty")
	} else if p.Prescription == "" {
		return errors.New("prescription cannot be empty")
	}
	return nil
}

type TimeRanges struct {
	Ranges []TimeRange
}

type UserType string

const (
	UserTypePatient    = "patient"
	UserTypeSpecialist = "specialist"
)

type AppointmentsRequest struct {
	ID       uuid.UUID `json:"id"` // patientID or specialistID
	UserType UserType  `json:"user_type"`
	Range    TimeRange `json:"range"`
}

type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Appointment struct {
	PatientID    uuid.UUID             `json:"patient_id"`
	SpecialistID uuid.UUID             `json:"specialist_id"`
	Speciality   models.Specialityenum `json:"speciality"`
	ScheduledAt  time.Time             `json:"scheduled_at"`
	Duration     Duration              `json:"duration"`

	PatientComment  string   `json:"patient_comment"`
	PatientSymptoms []string `json:"patient_symptoms"`
}

func (a Appointment) Invalidate() error {
	if a.PatientID == uuid.Nil {
		return errors.New("patient uuid cannot be nil")
	} else if a.SpecialistID == uuid.Nil {
		return errors.New("specialist uuid cannot be nil")
	} else if a.ScheduledAt.IsZero() {
		return errors.New("invalid scheduled time")
	} else if a.Duration.Seconds() == 0. {
		return errors.New("invalid duration")
	} else if a.PatientComment == "" {
		return errors.New("patient comment cannot be empty")
	} else if a.PatientSymptoms == nil || len(a.PatientSymptoms) == 0 {
		return errors.New("patient symptoms cannot be empty")
	}
	return nil
}

// Duration struct is used inside config.
type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
