package postgres

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/tools"
	"reflect"
)

type DataBaseParams struct {
	EnableLogging bool
}

type DataBase struct {
	db         *sql.DB
	dsn        string
	serverName string
	params     *DataBaseParams
}

func OpenDataBase(driver string, dsn string, params *DataBaseParams) (tools.DataBase, error) {
	return OpenDataBaseWithServerName(driver, "main", dsn, params)
}

func OpenDataBaseWithServerName(driver string, serverName string, dsn string, params *DataBaseParams) (tools.DataBase, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		logrus.Error("SQL", "connection to database failed", err)
		return nil, err
	}

	database := &DataBase{
		db:         db,
		dsn:        dsn,
		serverName: serverName,
		params:     params,
	}
	if params == nil {
		database.params = &DataBaseParams{EnableLogging: true}
	}

	if err := database.Ping(); err != nil {
		return nil, err
	}

	database.log("info", "database connection successful", dsn)
	return database, nil
}

func (d *DataBase) log(level string, query string, args ...any) {
	if d.params.EnableLogging {
		l := logrus.WithFields(logrus.Fields{"SQL": true, "Server": d.serverName})
		if level == "error" {
			l.Error(query, convertArgs(args))

		} else {
			l.Info(query, convertArgs(args))
		}
	}
}

func convertArgs(args ...any) []any {
	result := []any{}
	for _, val := range args {
		newVal := val
		if newVal != nil {
			if reflect.TypeOf(val) == reflect.TypeOf([]byte{}) {
				newVal = string(val.([]byte)[:])
			}
			if reflect.Array == reflect.TypeOf(newVal).Kind() || reflect.Slice == reflect.TypeOf(newVal).Kind() {
				newVal = convertArgs(newVal.([]any)...)
			}
		}
		result = append(result, newVal)
	}
	return result
}

func (d *DataBase) GetDB() *sql.DB {
	return d.db
}

func (d *DataBase) Exec(query string, args ...any) (sql.Result, error) {
	d.log("info", query, args...)
	result, err := d.db.Exec(query, args...)
	if err != nil {
		d.log("error", query, err)
	}
	return result, err
}

func (d *DataBase) Query(query string, args ...any) (*sql.Rows, error) {
	d.log("info", query, args...)
	result, err := d.db.Query(query, args...)
	if err != nil {
		d.log("error", query, err)
	}
	return result, err
}

func (d *DataBase) QueryRow(query string, args ...any) *sql.Row {
	d.log("info", query, args...)
	return d.db.QueryRow(query, args...)
}

func (d *DataBase) Ping() error {
	err := d.db.Ping()
	if err != nil {
		d.log("error", "ping database", err)
	} else {
		d.log("info", "ping database is ok")
	}
	return err
}

func (d *DataBase) Close() error {
	err := d.db.Close()
	if err != nil {
		d.log("error", "close database", err)
	} else {
		d.log("info", "connection to database is closed", d.dsn)
	}
	return err
}

func (d *DataBase) GetStats() string {
	return "unimplemented"
}
