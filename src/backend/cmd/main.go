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
	"github.com/go-chi/cors"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)

	sess := db.Init(log)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }

	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("drug", "drug", ""), configForNeo4j40)
	if err != nil {
		panic(err)
	}

	// For multidatabase support, set sessionConfig.DatabaseName to requested database
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "neo4j"}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		panic(err)
	}
	//fmt.Println(session.Run("call gds.graph.create('drug', ['Drug', 'Sub'], '*')", nil))

	cli := handler.NewHandler(sess, log, session)
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
		sess, log, 15*time.Second, p, mailCli,
	).Start()

	log.Info("running")
	if err := http.ListenAndServe(":8081", r); err != nil {
		panic(err)
	}
}
