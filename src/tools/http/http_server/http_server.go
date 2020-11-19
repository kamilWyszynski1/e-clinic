package http_server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type ServerConf struct {
	//Timeout for each request
	Timeout time.Duration
	//Health check function that will be mounted on GET /healthcheck (default return 200 handler)
	//Can be nil
	HeathCheck  http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
}

func HealthCheck(_ http.ResponseWriter, request *http.Request) {
	_, _ = ioutil.ReadAll(request.Body)
	_ = request.Body.Close()

	return
}

var DefaultConf = ServerConf{
	Timeout:    10 * time.Second,
	HeathCheck: HealthCheck,
	// maybe add middleware.RealIP, to default middleware?
}

//Provided config must include all fields (best way => copy default and change desired fields)
func New(log logrus.FieldLogger, config ...*ServerConf) (*chi.Mux, error) {
	conf := &DefaultConf
	if len(config) > 0 {
		conf = config[0]
	}

	r := chi.NewRouter()
	if conf.Timeout == 0 {
		return nil, fmt.Errorf("timeout must be greater than 0")
	}
	r.Use(middleware.Timeout(conf.Timeout), GetRecoverer(log))
	if len(conf.Middlewares) > 0 {
		r.Use(conf.Middlewares...)
	}

	if conf.HeathCheck != nil {
		r.Get("/healthcheck", conf.HeathCheck)
	}

	return r, nil
}

var ServerNotUp = errors.New("server is not working")

type HTTPServer struct {
	Logger    logrus.FieldLogger
	Server    *http.Server
	Mux       *chi.Mux
	Port      int
	Name      string
	WaitGroup *sync.WaitGroup
	// if no wait group was passed then wait for termination
	noWaitGroup bool
}

func NewServer(log logrus.FieldLogger, name string, port int, wg *sync.WaitGroup, configs ...*ServerConf) (*HTTPServer, error) {
	mux, err := New(log, configs...)
	if err != nil {
		return nil, fmt.Errorf("cannot create %s mux, due: %w", name, err)
	}
	wg2 := wg
	if wg == nil {
		wg2 = &sync.WaitGroup{}
	}
	return &HTTPServer{
		Logger: log,
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		Mux:         mux,
		Port:        port,
		Name:        name,
		WaitGroup:   wg2,
		noWaitGroup: wg == nil,
	}, nil
}

// StartServerRoutine just starts http server
func (s *HTTPServer) StartServerRoutine() {
	go func() {
		err := s.StartServer()
		if err != nil {
			s.Logger.WithError(err).Fatal("ServerRoutine failed!")
		}
	}()
}

// StartServer just starts http server
func (s *HTTPServer) StartServer() error {
	s.Logger.Infof("Start %s server at: %d", s.Name, s.Port)
	if s.WaitGroup != nil {
		s.WaitGroup.Add(1)
		defer s.WaitGroup.Done()
	}
	if err := s.Server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.Logger.Infof("%s server closed", s.Name)
			return nil
		} else {
			return fmt.Errorf("failed to listen %s server, due: %w", s.Name, err)
		}
	}
	return nil
}

func (s *HTTPServer) CloseServer(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("cannot shutdown server, due: %w", err)
	}
	if s.noWaitGroup {
		s.WaitGroup.Wait()
	}
	return nil
}
