package handler

import (
	"e-clinic/src/backend/clinic"

	"github.com/gocraft/dbr"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	db  *dbr.Session
	log logrus.FieldLogger
}

func NewHandler(db *dbr.Session, log logrus.FieldLogger) clinic.Assistant {
	return &Handler{db: db, log: log}
}
