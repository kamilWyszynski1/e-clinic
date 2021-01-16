package handler

import (
	"e-clinic/src/backend/clinic"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/gocraft/dbr"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	db     *dbr.Session
	log    logrus.FieldLogger
	neoCli neo4j.Session
	now    func() time.Time
}

func NewHandler(db *dbr.Session, log logrus.FieldLogger, neoCli neo4j.Session) clinic.Assistant {
	return &Handler{db: db, log: log, now: time.Now, neoCli: neoCli}
}
