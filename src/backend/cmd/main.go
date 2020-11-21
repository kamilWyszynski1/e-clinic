package main

import (
	"e-clinic/src/backend/clinic"
	"e-clinic/src/backend/clinic/handler"
	"e-clinic/src/backend/db"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.StandardLogger()

	sess := db.Init(log)

	r := chi.NewRouter()
	cli := handler.NewHandler(sess, log)
	clinic.RegisterAssistant(cli, r, log)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
