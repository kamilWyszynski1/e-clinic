package migration

// migrator This package handles migration in PostgreSQL database.
// Unfortunately go-pg/migrations do not support running operations from different directories.
// So usage os this directory is:
//
// 1) Create migrator.
// 2) Run go-pg/migrations using DB connection from Migrator and apply flag args.
// 3) Handle its response by migrator.
//
// {
//   ...
//   mg := migrator.CreateMigrator()
//   mg.HandleResponse(migrations.Run(mg.DbConnection, flag.Args()...))
//   ...
// }



import (
"flag"
"fmt"
"io/ioutil"
"os"
"strings"

"github.com/FatNinjas/SkyLane/src/tools/build_functions"
"github.com/FatNinjas/SkyLane/src/tools/database/dbconnection"
_ "github.com/FatNinjas/SkyLane/src/tools/filehandlers/importall"
"github.com/datainq/sdfmt"
"github.com/go-pg/migrations"
"github.com/go-pg/pg"
"github.com/sirupsen/logrus"
)

const usageText = `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
Usage:
  go run *.go <command> [args]
`

type Migrator struct {
	logger       logrus.FieldLogger
	DbConnection *pg.DB
	dir          string
}


func (m *Migrator) Run() {
	m.logger.Infof("migration dir: %s", m.dir)
	col := migrations.NewCollection()
	if err := col.DiscoverSQLMigrations(m.dir); err != nil {
		m.logger.WithError(err).Fatal("cannot discover migrations in dir")
	}
	m.logger.Info("listing migrations")
	for _, v := range col.Migrations() {
		m.logger.Infof("migration: %d (tx: %t)", v.Version, v.UpTx)
	}
	if fileList, err := ioutil.ReadDir(m.dir); err != nil {
		m.logger.WithError(err).Error("cannot list dir")
	} else {
		buf := &strings.Builder{}
		buf.WriteString("migration dir file list: \n")
		for _, l := range fileList {
			buf.WriteString(" * ")
			if l.IsDir() {
				buf.WriteString("(d) ")
			}
			buf.WriteString(l.Name())
			buf.WriteRune('\n')
		}
		m.logger.Info(buf.String())
	}
	m.logger.Info("args: ", flag.Args())
	m.HandleResponse(col.Run(m.DbConnection, flag.Args()...))
}

func (m *Migrator) HandleResponse(oldVersion, newVersion int64, err error) {
	statusCode := 0
	if err != nil {
		m.logger.Errorf(err.Error())
		statusCode = 1
	} else if newVersion != oldVersion {
		m.logger.Infof("migrated from version %d to %d", oldVersion, newVersion)
	} else {
		m.logger.Infof("version is %d", oldVersion)
	}
	logrus.Exit(statusCode)
}
