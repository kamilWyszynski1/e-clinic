package db

import (
	"os"
	"strings"

	"github.com/gocraft/dbr"
	"github.com/sirupsen/logrus"
)

func Init(logger logrus.FieldLogger) *dbr.Session {
	// INITIALIZE DB
	conn, err := dbr.Open("postgres", PostgreStringFromEnv(logger), nil)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to db")
	}
	if err := conn.Ping(); err != nil {
		logger.WithError(err).Fatal("could not ping db")
	} else {
		logger.Info("Connected to database!")
	}
	sess := conn.NewSession(nil)
	return sess
}

// This function creates PostreSQL connection string from environmental variables.
// It connects you to given database with given user.
// Desired variables are:
// * DB_HOST - http address of db instance
// * DB_PORT - port used by the db
// * DB_DATABASE - database that you want to connect
// * DB_USER - username to log
// * DB_PASSWORD - password of given above db user
// * DB_SSL_MODE - (optional) indicates mode of connection to db, default value: 'disable'
// * DB_SSL_CERT - (optional) path to ssl client certificate
// * DB_SSL_KEY - (optional) path to ssl user private key
// * DB_SSL_ROOT_CERT - (optional) path to ssl trusted root certificate
func PostgreStringFromEnv(log logrus.FieldLogger) string {
	log = log.WithFields(logrus.Fields{
		"tool":     "dbconnection",
		"database": "psql",
		"method":   "string",
	})
	connParams := map[string]string{
		"connect_timeout": "10", // Without it, the binary will hang long time waiting for a database connection.
	}
	getAndLog(log, connParams, "DB_HOST", "host")
	getAndLog(log, connParams, "DB_PORT", "port")
	getAndLog(log, connParams, "DB_DATABASE", "dbname")
	getAndLog(log, connParams, "DB_USER", "user")
	getAndLog(log, connParams, "DB_PASSWORD", "password")
	sslMode, _ := getAndLog(log, connParams, "DB_SSL_MODE", "sslmode")
	if sslMode == "allow" || sslMode == "prefer" || sslMode == "require" ||
		sslMode == "verify-ca" || sslMode == "verify-full" {

		getAndLog(log, connParams, "DB_SSL_CERT", "sslcert")
		getAndLog(log, connParams, "DB_SSL_KEY", "sslkey")
		getAndLog(log, connParams, "DB_SSL_ROOT_CERT", "sslrootcert")
	} else {
		connParams["sslmode"] = "disable"
	}
	builder := &strings.Builder{}
	params := &strings.Builder{}
	for k, v := range connParams {
		builder.WriteString(k)
		builder.WriteRune('=')
		builder.WriteString(v)
		builder.WriteRune(' ')

		params.WriteString(k)
		params.WriteRune(',')
	}
	connStr := builder.String()
	log.Infof(
		"Database Connection {host: %s, port: %s, user: '%s', pass: '%s', dbName: '%s', sslMode: '%s'}",
		connParams["host"], connParams["port"], connParams["user"], connParams["password"], connParams["dbname"], connParams["sslmode"],
	)
	return connStr
}

func getAndLog(log logrus.FieldLogger, m map[string]string, name, arg string) (string, bool) {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Warnf("ENV: %q missing", name)
	} else if value == "" {
		log.Warnf("ENV: %q present, but value is empty", name)
	} else {
		//log.Infof("ENV: %q present and value is %q", name, value)
		log.Debugf("ENV: %q present and value is %q", name, value)
		m[arg] = value
	}
	return value, ok
}
