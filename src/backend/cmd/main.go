package main

import (
	"e-clinic/src/backend/clinic"
	"e-clinic/src/backend/clinic/handler"
	"e-clinic/src/backend/clinic/mailing"
	"e-clinic/src/backend/clinic/payment"
	"e-clinic/src/backend/db"
	payugo "e-clinic/src/backend/payu"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.StandardLogger()

	sess := db.Init(log)

	r := chi.NewRouter()
	cli := handler.NewHandler(sess, log)
	clinic.RegisterAssistant(cli, r, log)

	// PAYMENT
	p, err := payugo.NewClient(
		http.DefaultClient, "https://secure.snd.payu.com",
		payugo.MerchantConfig{
			ClientID:     "398268",
			ClientSecret: "880487191465ca9418fafcd9c0a019e6",
			PosID:        "398268",
		})
	if err != nil {
		panic(err)
	}
	if err := p.Authorize(); err != nil {
		panic(err)
	}

	// MAILING
	mailCli := mailing.NewClient(mailjet.NewMailjetClient("fb71068ebf8203243a86c64e951f7778", "3450a83ffd0cf668ded207e42f46830b"))

	payment.NewWatcher(
		sess, log, 10*time.Second, p, mailCli,
	).Start()

	log.Info("running")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
