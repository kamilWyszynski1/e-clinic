package clinic

import (
	"e-clinic/src/models"
	"encoding/json"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Caller interface {
	// Call method calls specific person to inform him about status change in appointment
	Call(id uuid.UUID, appointmentID uuid.UUID, action Action) error
}

type Action string

const (
	ActionAppointmentScheduled = "scheduled"
	ActionAppointmentAccepted  = "accepted"
)

type Assistant interface {
	// GetAppointments returns appointments in given time range
	GetAppointments(t *TimeRange) (*AppointmentList, int, error)
}

type AppointmentList struct {
	Appointments []models.Appointment `json:"appointments"`
	Len          int                  `json:"len"`
}

type SpecjalistAssistant interface {
	GetFreeTime(id uuid.UUID, timeRange *TimeRange) (*TimeRanges, int, error)
}

type TimeRanges struct {
	Ranges []TimeRange
}

type PatientAssistant interface {
	// CreateAppointment checks if appointment is valid and schedules it
	CreateAppointment(a *Appointment) (int, error)
	GetAppointments(t *TimeRange) (*AppointmentList, int, error)
}

type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Appointment struct {
	PatientID    uuid.UUID `json:"patient_id"`
	SpecialistID uuid.UUID `json:"specialist_id"`
	ScheduledAt  time.Time `json:"scheduled_at"`
	Duration     Duration  `json:"duration"`

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
