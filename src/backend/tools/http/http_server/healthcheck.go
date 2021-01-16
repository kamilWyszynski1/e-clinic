package http_server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func GetDefaultLocalHeathCheckAddress(port int) string {
	return fmt.Sprintf("http://127.0.0.1:%d/healthcheck", port)
}

func NewMultipleMuxHealthCheckWithTurnOff(logger logrus.FieldLogger, serverAddrToCheck []string, isApiOn *bool) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if *isApiOn {
			_, _ = ioutil.ReadAll(request.Body)
			_ = request.Body.Close()
			for _, addr := range serverAddrToCheck {
				// check server
				resp, err := http.DefaultClient.Get(addr)
				if err != nil {
					logger.WithError(err).Errorf("cannot request %s server!", addr)
					http.Error(writer, "", http.StatusInternalServerError)
					return
				}
				// check if health check was successful was ok
				if resp.StatusCode != http.StatusOK {
					logger.WithField("statusCode", resp.StatusCode).Errorf("server %s is not ok", addr)
					http.Error(writer, "", resp.StatusCode)
					return
				}
				// close response
				_, _ = ioutil.ReadAll(resp.Body)
				_ = resp.Body.Close()
			}
			writer.WriteHeader(http.StatusOK)
		} else {
			logger.Warn("API is turned off")
			http.Error(writer, "server is off", http.StatusServiceUnavailable)
		}
	}
}
