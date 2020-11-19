package specjalist

import (
	"e-clinic/src/clinic"
	"net/http"
	"time"

	"github.com/gocraft/dbr"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	db  *dbr.Session
	log logrus.FieldLogger
}

func NewSpecjalistHandler(db *dbr.Session, log logrus.FieldLogger) clinic.SpecjalistAssistant {
	return &Handler{db: db, log: log}
}

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
