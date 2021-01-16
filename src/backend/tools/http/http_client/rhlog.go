package http_client

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

type rhlog struct {
	log logrus.FieldLogger
}

func (r rhlog) Error(msg string, args ...interface{}) {
	r.log.WithField("msg", msg).Warning(args)
}

func (r rhlog) Info(msg string, args ...interface{}) {
	r.log.WithField("msg", msg).Info(args)
}

func (r rhlog) Debug(msg string, args ...interface{}) {
	r.log.WithField("msg", msg).Debug(args)
}

func (r rhlog) Warn(msg string, args ...interface{}) {
	r.log.WithField("msg", msg).Warning(args)
}

func NewRHLog(log logrus.FieldLogger) retryablehttp.LeveledLogger {
	return &rhlog{log: log}
}
