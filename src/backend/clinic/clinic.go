package clinic

import (
	"e-clinic/src/backend/models"
	"encoding/json"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Assistant interface {
	// GetSpecialistFreeTime returns specialist's free time
	GetSpecialistFreeTime(id uuid.UUID, timeRange *TimeRange) (*TimeRanges, int, error)
	// MakePrescription creates prescription
	MakePrescription(p *Prescription) (int, error)
	// AcceptAppointment specialist accepts appointment
	AcceptAppointment(aID uuid.UUID) (int, error)
	// RejectAppointment specialist rejects appointment
	RejectAppointment(aID uuid.UUID) (int, error)

	// CreateAppointment checks if appointment is valid and schedules it
	CreateAppointment(a *Appointment) (*models.Appointment, int, error)
	// GetAppointments returns user's appointments in given time range
	GetAppointments(ar *AppointmentsRequest) (*AppointmentList, int, error)
	// GetAppointment returns appointment details
	GetAppointment(aID uuid.UUID) (*AppointmentInfo, int, error)
	// GetSpecialists returns specialists
	GetSpecialists() (*SpecialistList, int, error)

	// GetDrugs returns drug list with pagination
	GetDrugs(prefix string, offset int, limit int) (*Drugs, int, error)
	// GetDrug returns drug's info
	GetDrug(drugID int) (*DrugWithSubstances, int, error)
	// GetReplacement returns drug's replacements
	GetReplacement(drugID int, minSimilarity float64) (*Drugs, int, error)
}

type SpecialistList struct {
	Specialists []SpecialistWithFee `json:"specialists"`
}

type SpecialistWithFee struct {
	SpecialistID uuid.UUID             `json:"specialist_id,omitempty"` // id
	FeeId        uuid.UUID             `json:"fee_id"`
	Name         string                `json:"name,omitempty"`    // name
	Surname      string                `json:"surname,omitempty"` // surname
	Speciality   models.Specialityenum `json:"speciality"`
	FeePer30Min  float64               `json:"fee_per_30_min"`
}

type AppointmentInfo struct {
	Appointment  *models.Appointment     `json:"appointment"`
	Form         *models.AppointmentForm `json:"form"`
	Prescription *Prescription           `json:"prescription"`
}

type DrugWithSubstances struct {
	Drug       *models.Drug        `json:"drug"`
	Substances []*models.Substance `json:"substances"`
}

type Substance struct {
	Name string `json:"name,omitempty"` // name
}

type Drugs struct {
	Drugs []*models.Drug `json:"drugs"`
	Len   int            `json:"len"`
}

type AppointmentList struct {
	Appointments []*AppointmentInfo `json:"appointments"`
	Len          int                `json:"len"`
}

type Prescription struct {
	AppointmentID uuid.UUID  `json:"appointment_id,omitempty"`
	SpecialistID  uuid.UUID  `json:"specialist_id,omitempty"`
	Comment       string     `json:"comment"`
	Drugs         []DrugDose `json:"drugs"`
}

type DrugDose struct {
	Drug   int    `json:"drug"` // listed drugs
	Dosing string `json:"dosing"`
}

func (p Prescription) Invalidate() error {
	if p.Comment == "" {
		return errors.New("comment cannot be empty")
	} else if len(p.Drugs) == 0 {
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

func (a Appointment) Validate(now time.Time) error {
	if a.PatientID == uuid.Nil {
		return errors.New("patient uuid cannot be nil")
	} else if a.SpecialistID == uuid.Nil {
		return errors.New("specialist uuid cannot be nil")
	} else if a.ScheduledAt.IsZero() {
		return errors.New("invalid scheduled time")
	} else if a.Duration.Seconds() == 0. {
		return errors.New("invalid duration")
	} else if a.ScheduledAt.Before(now) {
		return errors.New("appointment must be in future")
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
