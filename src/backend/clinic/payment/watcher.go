package payment

import (
	"database/sql"
	"e-clinic/src/backend/clinic/handler"
	"e-clinic/src/backend/clinic/mailing"
	"e-clinic/src/backend/models"
	payugo "e-clinic/src/backend/payu"
	"fmt"
	"strconv"
	"time"

	"github.com/gocraft/dbr"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type newPayment struct {
	AppointmentID uuid.UUID
	Price         float64
	Name, Email   string
}

type channels struct {
	AcceptedCh chan *newPayment
	ToBePaidCh chan string
}

type Watcher struct {
	db         *dbr.Session
	log        logrus.FieldLogger
	sleep      time.Duration
	ch         channels
	paymentCli *payugo.Client
	mailCli    mailing.Mailer
}

func NewWatcher(db *dbr.Session, log logrus.FieldLogger, sleep time.Duration, paymentCli *payugo.Client, mailCli mailing.Mailer) *Watcher {
	return &Watcher{
		db:    db,
		log:   log,
		sleep: sleep,
		ch: channels{
			AcceptedCh: make(chan *newPayment, 20),
			ToBePaidCh: make(chan string),
		},
		paymentCli: paymentCli,
		mailCli:    mailCli,
	}
}

// Start starts watching over payment
// watcher queries appointment to create payment
// and watches over created transactions
func (w Watcher) Start() {
	go func() {
		for {
			select {
			case <-time.After(w.sleep):
				w.log.Info("querying")
				w.queryPayments(w.log.WithField("method", "queryPayments"))
			}
		}
	}()

	go w.watchOverPayments()
}

const queryAppointmentPayments = `SELECT a.id as appointment_id, a.state, a.duration as duration, p.order_id, sf.fee_per_30_min as fee, p2.email, p2.name FROM appointment a
JOIN specialist_fee sf on a.specialist_fee = sf.id
LEFT JOIN payment p on a.id = p.appointment
JOIN patient p2 on a.patient = p2.id
WHERE a.state IN ('ACCEPTED', 'TO_BE_PAID')`

type appointmentPayment struct {
	AppointmentID uuid.UUID
	State         models.Apoitntmentstateenum
	Duration      int            `json:"duration"`
	Fee           float64        `json:"fee"` // fee_per_30_min
	OrderID       sql.NullString `json:"order_id"`
	Email         string
	Name          string
}

func (w Watcher) queryPayments(log *logrus.Entry) {
	var payments []appointmentPayment
	_, err := w.db.SelectBySql(queryAppointmentPayments).Load(&payments)
	if err != nil {
		log.WithError(err).Error("failed to query payments")
		return
	}

	for _, p := range payments {
		fmt.Printf("%+v\n", p)
		switch p.State {
		case models.ApoitntmentstateenumAccepted: // create new payment
			w.ch.AcceptedCh <- &newPayment{
				AppointmentID: p.AppointmentID,
				Price:         p.Fee * float64(p.Duration) / 60. / 60.,
				Name:          p.Email,
				Email:         p.Name,
			}
		case models.ApoitntmentstateenumToBePaid: // check if payment passed
			if p.OrderID.Valid && p.OrderID.String != "" {
				w.ch.ToBePaidCh <- p.OrderID.String
			} else {
				w.log.Error("empty orderID when appointment in ToBePaid state")
			}
		default:
			log.Errorf("invalid appointment state: %s", p.State)
			return
		}
	}
}

func (w Watcher) watchOverPayments() {
	log := w.log.WithField("method", "watchOverPayments")
	for {
		select {
		case a := <-w.ch.AcceptedCh:
			log.Info("creating new payment")
			fmt.Println(a.Price)
			order, err := w.paymentCli.OrderCreateRequest(&payugo.Order{
				Description:  fmt.Sprintf("Payment for appointment: %s", a.AppointmentID),
				CurrencyCode: payugo.CurrencyCodePLN,
				TotalAmount:  strconv.Itoa(int(a.Price)),
				ContinueURL:  "",
				ExtOrderID:   "testqweqwe123123reterter",
				CustomerIP:   "127.0.0.1",
				Buyer: payugo.Buyer{
					Email:     "kamil.wyszynski.97@gmail.com",
					Phone:     "",
					FirstName: "Kamil",
					LastName:  "Wyszyński",
					Language:  payugo.LanguagePL,
				},
				Products: []payugo.Product{
					{
						Name:      fmt.Sprintf("Appointment: %s", a.AppointmentID),
						UnitPrice: strconv.Itoa(int(a.Price) * 1000),
						Quantity:  "1",
					},
				},
			})
			if err != nil {
				log.WithError(err).Error("failed to create new payment")
				continue
			}

			_, err = w.mailCli.SendEmail(
				a.Name, a.Email, "E-clinic payment",
				fmt.Sprintf("your payment, click here: %s", order.RedirectURI))
			if err != nil {
				log.WithError(err).Error("failed to send email")
				continue
			}

			p := models.Payment{
				ID:          uuid.NewV4(),
				Appointment: a.AppointmentID,
				Price:       a.Price,
				OrderID:     order.OrderID,
				Status:      order.Status.StatusDesc + ";" + order.Status.StatusCode.String(),
			}
			if err := p.Insert(w.db); err != nil {
				log.WithError(err).Error("failed to insert new payment")
				continue
			}

			if err := handler.ChangeAppointmentStatus(w.db, a.AppointmentID, models.ApoitntmentstateenumToBePaid); err != nil {
				log.WithError(err).Error("failed to change appointment state")
				continue
			}

			log.Info("created new payment")
		}
	}
}
