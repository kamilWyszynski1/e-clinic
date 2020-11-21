package clinic

import (
	"e-clinic/src/backend/models"
	"e-clinic/src/backend/tools/http/http_client"
	"e-clinic/src/backend/tools/http/http_server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const AssistantRoutePattern = "/api/v1/Assistant"

func RegisterAssistant(instance Assistant, r *chi.Mux, log logrus.FieldLogger) {
	r.Route("/api/v1/Assistant", func(r chi.Router) {
		r.Post("/GetFreeTime", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument id
			idStr := request.URL.Query().Get("id")
			id, err := uuid.FromString(idStr)
			if err != nil {
				log.WithError(err).Error("can't parse uuid")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			// Reading argument timeRange
			rawBody, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.WithError(err).Error("can't read body")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			var timeRange *TimeRange
			if len(rawBody) != 0 {
				timeRange = &TimeRange{}
				if err := json.Unmarshal(rawBody, timeRange); err != nil {
					log.WithError(err).Error("can't unmarshal body")
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			// Calling function GetFreeTime
			{
				response, responseCode, err := instance.GetFreeTime(id, timeRange)
				writer.WriteHeader(responseCode)
				if err != nil {
					if err := json.NewEncoder(writer).Encode(&struct{ Error string }{Error: err.Error()}); err != nil {
						fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
						log.WithError(err).Error("Could not marshal response error")
					}
					return
				}
				if response != nil {
					if err := json.NewEncoder(writer).Encode(response); err != nil {
						log.WithError(err).Error("Could not write response")
						writer.WriteHeader(http.StatusInternalServerError)
						if err := json.NewEncoder(writer).Encode(&struct{ Error string }{Error: err.Error()}); err != nil {
							fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
							log.WithError(err).Error("Could not marshal response error")
						}
						return
					}
				}
			}
		})
		r.Post("/MakePrescription", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument p
			rawBody, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.WithError(err).Error("can't read body")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			var p *Prescription
			if len(rawBody) != 0 {
				p = &Prescription{}
				if err := json.Unmarshal(rawBody, p); err != nil {
					log.WithError(err).Error("can't unmarshal body")
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			// Calling function MakePrescription
			{
				responseCode, err := instance.MakePrescription(p)
				writer.WriteHeader(responseCode)
				if err != nil {
					if err := json.NewEncoder(writer).Encode(&struct{ Error string }{Error: err.Error()}); err != nil {
						fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
						log.WithError(err).Error("Could not marshal response error")
					}
					return
				}
			}
		})
		r.Get("/AcceptAppointment", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument aID
			aIDStr := request.URL.Query().Get("aID")
			aID, err := uuid.FromString(aIDStr)
			if err != nil {
				log.WithError(err).Error("can't parse uuid")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			// Calling function AcceptAppointment
			{
				responseCode, err := instance.AcceptAppointment(aID)
				writer.WriteHeader(responseCode)
				if err != nil {
					if err := json.NewEncoder(writer).Encode(&struct{ Error string }{Error: err.Error()}); err != nil {
						fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
						log.WithError(err).Error("Could not marshal response error")
					}
					return
				}
			}
		})
		r.Get("/RejectAppointment", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument aID
			aIDStr := request.URL.Query().Get("aID")
			aID, err := uuid.FromString(aIDStr)
			if err != nil {
				log.WithError(err).Error("can't parse uuid")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			// Calling function RejectAppointment
			{
				responseCode, err := instance.RejectAppointment(aID)
				writer.WriteHeader(responseCode)
				if err != nil {
					if err := json.NewEncoder(writer).Encode(&struct{ Error string }{Error: err.Error()}); err != nil {
						fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
						log.WithError(err).Error("Could not marshal response error")
					}
					return
				}
			}
		})
		r.Post("/CreateAppointment", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument a
			rawBody, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.WithError(err).Error("can't read body")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			var a *Appointment
			if len(rawBody) != 0 {
				a = &Appointment{}
				if err := json.Unmarshal(rawBody, a); err != nil {
					log.WithError(err).Error("can't unmarshal body")
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			// Calling function CreateAppointment
			{
				response, responseCode, err := instance.CreateAppointment(a)
				writer.WriteHeader(responseCode)
				writer.Header().Add("Content-Type", "application/json")
				if err != nil {
					if err := json.NewEncoder(writer).Encode(&struct{ Error string }{Error: err.Error()}); err != nil {
						fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
						log.WithError(err).Error("Could not marshal response error")
					}
					return
				}
				if response != nil {
					if err := json.NewEncoder(writer).Encode(response); err != nil {
						log.WithError(err).Error("Could not write response")
						writer.WriteHeader(http.StatusInternalServerError)
						if err := json.NewEncoder(writer).Encode(&struct{ Error string }{Error: err.Error()}); err != nil {
							fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
							log.WithError(err).Error("Could not marshal response error")
						}
						return
					}
				}
			}
		})
		r.Post("/GetAppointments", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument ar
			rawBody, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.WithError(err).Error("can't read body")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			var ar *AppointmentsRequest
			if len(rawBody) != 0 {
				ar = &AppointmentsRequest{}
				if err := json.Unmarshal(rawBody, ar); err != nil {
					log.WithError(err).Error("can't unmarshal body")
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			// Calling function GetAppointments
			{
				response, responseCode, err := instance.GetAppointments(ar)
				writer.WriteHeader(responseCode)
				if err != nil {
					if err := json.NewEncoder(writer).Encode(&struct{ Error string }{Error: err.Error()}); err != nil {
						fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
						log.WithError(err).Error("Could not marshal response error")
					}
					return
				}
				if response != nil {
					if err := json.NewEncoder(writer).Encode(response); err != nil {
						log.WithError(err).Error("Could not write response")
						writer.WriteHeader(http.StatusInternalServerError)
						if err := json.NewEncoder(writer).Encode(&struct{ Error string }{Error: err.Error()}); err != nil {
							fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
							log.WithError(err).Error("Could not marshal response error")
						}
						return
					}
				}
			}
		})
		r.Get("/isUp", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		})
	})
}

type AssistantRestClient struct {
	address string
	client  http_client.DefaultHTTPClient
	errs    []error
}

func NewAssistantRestClient(address string, client http_client.DefaultHTTPClient) (*AssistantRestClient, error) {
	packageErrors := []error{}

	restClient := &AssistantRestClient{
		address: address,
		client:  client,
		errs:    packageErrors,
	}

	// restClient is returned to avoid panics, some services may work without working rest client
	if code := restClient.isUp(); code != http.StatusOK {
		return restClient, http_server.ServerNotUp
	}
	return restClient, nil
}

func NewAssistantRestClientErrWrap(address string, client http_client.DefaultHTTPClient, log logrus.FieldLogger) *AssistantRestClient {
	restClient, err := NewAssistantRestClient(address, client)
	if err != nil {
		log.WithError(err).Warn("failed to create AssistantRestClient")
	}
	return restClient
}

func (serviceInstance *AssistantRestClient) wrapError(errMsg interface{}) error {
	for _, er := range serviceInstance.errs {
		if strings.EqualFold(errMsg.(string), er.Error()) {
			return er
		}
	}
	return fmt.Errorf("%v", errMsg)
}

func (serviceInstance *AssistantRestClient) GetFreeTime(id uuid.UUID, timeRange *TimeRange) (*TimeRanges, int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/GetFreeTime")
	if err != nil {
		return nil, -1, err
	}
	query := u.Query()
	query.Set("id", id.String())

	u.RawQuery = query.Encode()

	data, err := json.Marshal(timeRange)
	if err != nil {
		return nil, -1, err
	}
	resp := serviceInstance.client.Post(u.String(), data)
	if resp.Err != nil {
		return nil, -1, resp.Err
	}

	if len(resp.Body) > 0 {
		respErr := map[string]interface{}{}
		if err := json.Unmarshal(resp.Body, &respErr); err != nil {
			return nil, resp.StatusCode, err
		}
		if errorMessage, found := respErr["Error"]; found && len(respErr) == 1 {
			return nil, resp.StatusCode, serviceInstance.wrapError(errorMessage)
		}
	}
	respJson := &TimeRanges{}
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, respJson); err != nil {
			return nil, resp.StatusCode, err
		}
	}
	return respJson, resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) MakePrescription(p *Prescription) (int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/MakePrescription")
	if err != nil {
		return -1, err
	}

	data, err := json.Marshal(p)
	if err != nil {
		return -1, err
	}
	resp := serviceInstance.client.Post(u.String(), data)
	if resp.Err != nil {
		return -1, resp.Err
	}

	if len(resp.Body) > 0 {
		respErr := map[string]interface{}{}
		if err := json.Unmarshal(resp.Body, &respErr); err != nil {
			return resp.StatusCode, err
		}
		if errorMessage, found := respErr["Error"]; found && len(respErr) == 1 {
			return resp.StatusCode, serviceInstance.wrapError(errorMessage)
		}
	}
	return resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) AcceptAppointment(aID uuid.UUID) (int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/AcceptAppointment")
	if err != nil {
		return -1, err
	}
	query := u.Query()
	query.Set("aID", aID.String())

	u.RawQuery = query.Encode()

	resp := serviceInstance.client.Get(u.String())
	if resp.Err != nil {
		return -1, resp.Err
	}

	if len(resp.Body) > 0 {
		respErr := map[string]interface{}{}
		if err := json.Unmarshal(resp.Body, &respErr); err != nil {
			return resp.StatusCode, err
		}
		if errorMessage, found := respErr["Error"]; found && len(respErr) == 1 {
			return resp.StatusCode, serviceInstance.wrapError(errorMessage)
		}
	}
	return resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) RejectAppointment(aID uuid.UUID) (int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/RejectAppointment")
	if err != nil {
		return -1, err
	}
	query := u.Query()
	query.Set("aID", aID.String())

	u.RawQuery = query.Encode()

	resp := serviceInstance.client.Get(u.String())
	if resp.Err != nil {
		return -1, resp.Err
	}

	if len(resp.Body) > 0 {
		respErr := map[string]interface{}{}
		if err := json.Unmarshal(resp.Body, &respErr); err != nil {
			return resp.StatusCode, err
		}
		if errorMessage, found := respErr["Error"]; found && len(respErr) == 1 {
			return resp.StatusCode, serviceInstance.wrapError(errorMessage)
		}
	}
	return resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) CreateAppointment(a *Appointment) (*models.Appointment, int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/CreateAppointment")
	if err != nil {
		return nil, -1, err
	}

	data, err := json.Marshal(a)
	if err != nil {
		return nil, -1, err
	}
	resp := serviceInstance.client.Post(u.String(), data)
	if resp.Err != nil {
		return nil, -1, resp.Err
	}

	if len(resp.Body) > 0 {
		respErr := map[string]interface{}{}
		if err := json.Unmarshal(resp.Body, &respErr); err != nil {
			return nil, resp.StatusCode, err
		}
		if errorMessage, found := respErr["Error"]; found && len(respErr) == 1 {
			return nil, resp.StatusCode, serviceInstance.wrapError(errorMessage)
		}
	}
	respJson := &models.Appointment{}
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, respJson); err != nil {
			return nil, resp.StatusCode, err
		}
	}
	return respJson, resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) GetAppointments(ar *AppointmentsRequest) (*AppointmentList, int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/GetAppointments")
	if err != nil {
		return nil, -1, err
	}

	data, err := json.Marshal(ar)
	if err != nil {
		return nil, -1, err
	}
	resp := serviceInstance.client.Post(u.String(), data)
	if resp.Err != nil {
		return nil, -1, resp.Err
	}

	if len(resp.Body) > 0 {
		respErr := map[string]interface{}{}
		if err := json.Unmarshal(resp.Body, &respErr); err != nil {
			return nil, resp.StatusCode, err
		}
		if errorMessage, found := respErr["Error"]; found && len(respErr) == 1 {
			return nil, resp.StatusCode, serviceInstance.wrapError(errorMessage)
		}
	}
	respJson := &AppointmentList{}
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, respJson); err != nil {
			return nil, resp.StatusCode, err
		}
	}
	return respJson, resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) isUp() int {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/isUp")
	if err != nil {
		return -1
	}

	resp := serviceInstance.client.Get(u.String())
	if resp.Err != nil {
		return -1
	}

	return resp.StatusCode
}
