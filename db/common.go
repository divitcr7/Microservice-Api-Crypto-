package db

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/streamdp/ccd/config"
	"github.com/streamdp/ccd/db/mysql"
	"github.com/streamdp/ccd/db/postgres"
	"github.com/streamdp/ccd/domain"
)

// Session interface makes it possible to expand the list of session storages
type Session interface {
	AddTask(n string, i int64) (err error)
	UpdateTask(n string, i int64) (err error)
	RemoveTask(n string) (err error)
	GetSession() (map[string]int64, error)
}

// Database interface makes it possible to expand the list of data storages
type Database interface {
	Insert(data *domain.Data) (result sql.Result, err error)
	GetLast(from string, to string) (result *domain.Data, err error)
	DataPipe() chan *domain.Data

	AddSymbol(s string, u string) (result sql.Result, err error)
	UpdateSymbol(s string, u string) (result sql.Result, err error)
	RemoveSymbol(s string) (result sql.Result, err error)
	Symbols() (symbols []*domain.Symbol, err error)

	AddTask(n string, i int64) (result sql.Result, err error)
	UpdateTask(n string, i int64) (result sql.Result, err error)
	RemoveTask(n string) (result sql.Result, err error)
	GetSession() (tasks map[string]int64, err error)
}

func Connect(l *log.Logger) (d Database, err error) {
	var (
		driverName       = mysql.Mysql
		dataBaseUrl      = config.GetEnv("CCDC_DATABASEURL")
		connectionString string
	)
	if dataBaseUrl == "" {
		return nil, errors.New("please set OS environment \"CCDC_DATABASEURL\" with database connection string")
	}
	connectionParameters := strings.Split(dataBaseUrl, "://")
	if len(connectionParameters) == 2 {
		driverName, connectionString = connectionParameters[0], connectionParameters[1]
	}
	switch driverName {
	case postgres.Postgres:
		d, err = postgres.Connect(dataBaseUrl)
	case mysql.Mysql:
		fallthrough
	default:
		d, err = mysql.Connect(connectionString)
	}
	if err == nil {
		serve(d, l)
	}
	return
}

func serve(d Database, l *log.Logger) {
	go func() {
		for data := range d.DataPipe() {
			if _, err := d.Insert(data); err != nil {
				l.Println(err)
			}
		}
	}()
}
