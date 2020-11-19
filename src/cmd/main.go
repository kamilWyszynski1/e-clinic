package main

import (
	"e-clinic/src/clinic"
	"e-clinic/src/clinic/patient"
	"e-clinic/src/clinic/specjalist"
	"e-clinic/src/db"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.StandardLogger()

	sess := db.Init(log)

	r := chi.NewRouter()
	sCli := specjalist.NewSpecjalistHandler(sess, log)
	pCli := patient.NewPatientHandler(sess, log, sCli)

	clinic.RegisterPatientAssistant(pCli, r, log)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
