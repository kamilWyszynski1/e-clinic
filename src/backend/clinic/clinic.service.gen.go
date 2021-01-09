package clinic

// Code generated by rest_rpc. DO NOT EDIT.

import (
	"e-clinic/src/backend/models"
	"e-clinic/src/backend/tools/http/http_client"
	"e-clinic/src/backend/tools/http/http_server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const AssistantRoutePattern = "/api/v1/Assistant"

func RegisterAssistant(instance Assistant, r *chi.Mux, log logrus.FieldLogger) {
	r.Route("/api/v1/Assistant", func(r chi.Router) {
		r.Post("/GetSpecialistFreeTime", func(writer http.ResponseWriter, request *http.Request) {
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

			// Calling function GetSpecialistFreeTime
			{
				response, responseCode, err := instance.GetSpecialistFreeTime(id, timeRange)
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
		r.Get("/GetAppointment", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument aID
			aIDStr := request.URL.Query().Get("aID")
			aID, err := uuid.FromString(aIDStr)
			if err != nil {
				log.WithError(err).Error("can't parse uuid")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			// Calling function GetAppointment
			{
				response, responseCode, err := instance.GetAppointment(aID)
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
		r.Get("/GetSpecialists", func(writer http.ResponseWriter, request *http.Request) {
			// Calling function GetSpecialists
			{
				response, responseCode, err := instance.GetSpecialists()
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
		r.Get("/GetDrugs", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument prefix
			prefix := request.URL.Query().Get("prefix")
			// Reading argument offset
			offsetStr := request.URL.Query().Get("offset")
			offset, err := strconv.Atoi(offsetStr)
			if err != nil {
				log.WithError(err).Errorf("Could not parse offset parameter value. Expected integer got %s", offsetStr)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			// Reading argument limit
			limitStr := request.URL.Query().Get("limit")
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				log.WithError(err).Errorf("Could not parse limit parameter value. Expected integer got %s", limitStr)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			// Calling function GetDrugs
			{
				response, responseCode, err := instance.GetDrugs(prefix, offset, limit)
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
		r.Get("/GetDrug", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument drugID
			drugIDStr := request.URL.Query().Get("drugID")
			drugID, err := strconv.Atoi(drugIDStr)
			if err != nil {
				log.WithError(err).Errorf("Could not parse drugID parameter value. Expected integer got %s", drugIDStr)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			// Calling function GetDrug
			{
				response, responseCode, err := instance.GetDrug(drugID)
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
		r.Get("/GetReplacement", func(writer http.ResponseWriter, request *http.Request) {
			// Reading argument drugID
			drugIDStr := request.URL.Query().Get("drugID")
			drugID, err := strconv.Atoi(drugIDStr)
			if err != nil {
				log.WithError(err).Errorf("Could not parse drugID parameter value. Expected integer got %s", drugIDStr)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			// Reading argument minSimilarity
			minSimilarityStr := request.URL.Query().Get("minSimilarity")
			minSimilarity, err := strconv.ParseFloat(minSimilarityStr, 64)
			if err != nil {
				log.WithError(err).Errorf("Could not parse minSimilarity parameter value. Expected float got %s", minSimilarityStr)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			// Calling function GetReplacement
			{
				response, responseCode, err := instance.GetReplacement(drugID, minSimilarity)
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

func (serviceInstance *AssistantRestClient) GetSpecialistFreeTime(id uuid.UUID, timeRange *TimeRange) (*TimeRanges, int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/GetSpecialistFreeTime")
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

func (serviceInstance *AssistantRestClient) GetAppointment(aID uuid.UUID) (*AppointmentInfo, int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/GetAppointment")
	if err != nil {
		return nil, -1, err
	}
	query := u.Query()
	query.Set("aID", aID.String())

	u.RawQuery = query.Encode()

	resp := serviceInstance.client.Get(u.String())
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
	respJson := &AppointmentInfo{}
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, respJson); err != nil {
			return nil, resp.StatusCode, err
		}
	}
	return respJson, resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) GetSpecialists() (*SpecialistList, int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/GetSpecialists")
	if err != nil {
		return nil, -1, err
	}

	resp := serviceInstance.client.Get(u.String())
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
	respJson := &SpecialistList{}
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, respJson); err != nil {
			return nil, resp.StatusCode, err
		}
	}
	return respJson, resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) GetDrugs(prefix string, offset int, limit int) (*Drugs, int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/GetDrugs")
	if err != nil {
		return nil, -1, err
	}
	query := u.Query()
	query.Set("prefix", prefix)
	query.Set("offset", strconv.Itoa(offset))
	query.Set("limit", strconv.Itoa(limit))

	u.RawQuery = query.Encode()

	resp := serviceInstance.client.Get(u.String())
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
	respJson := &Drugs{}
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, respJson); err != nil {
			return nil, resp.StatusCode, err
		}
	}
	return respJson, resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) GetDrug(drugID int) (*DrugWithSubstances, int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/GetDrug")
	if err != nil {
		return nil, -1, err
	}
	query := u.Query()
	query.Set("drugID", strconv.Itoa(drugID))

	u.RawQuery = query.Encode()

	resp := serviceInstance.client.Get(u.String())
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
	respJson := &DrugWithSubstances{}
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, respJson); err != nil {
			return nil, resp.StatusCode, err
		}
	}
	return respJson, resp.StatusCode, nil
}

func (serviceInstance *AssistantRestClient) GetReplacement(drugID int, minSimilarity float64) (*Drugs, int, error) {
	u, err := url.Parse(serviceInstance.address + "/api/v1/Assistant/GetReplacement")
	if err != nil {
		return nil, -1, err
	}
	query := u.Query()
	query.Set("drugID", strconv.Itoa(drugID))
	query.Set("minSimilarity", strconv.FormatFloat(minSimilarity, 'f', -1, 64))

	u.RawQuery = query.Encode()

	resp := serviceInstance.client.Get(u.String())
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
	respJson := &Drugs{}
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
