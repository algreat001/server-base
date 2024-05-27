package tools

import (
	"errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	cacheTimeLife = 1 * time.Hour
	indexError    = errors.New("index not found")
	dataBaseError = errors.New("database not found")
	READ          = "read"
	WRITE         = "write"
)

type TypeQuery struct {
	typeQuery string
	cached    time.Time
}

type TypeQueryPool struct {
	typeQuery map[string]*TypeQuery
}

func NewTypeQueryPool() *TypeQueryPool {
	return &TypeQueryPool{
		typeQuery: make(map[string]*TypeQuery),
	}
}

func (tq *TypeQueryPool) GetTypeQuery(query string, db DataBase) string {
	funcName, err := getQueryFunction(query)
	if err != nil {
		logrus.Info("error get query function: ", err)
		return WRITE
	}

	if tq.isCached(funcName) {
		return tq.typeQuery[funcName].typeQuery
	}
	err = tq.updateTypeQuery(funcName, db)
	if err != nil {
		logrus.Info("error update type query: ", err)
		return WRITE
	}
	return tq.typeQuery[funcName].typeQuery
}

func (tq *TypeQueryPool) updateTypeQuery(funcName string, db DataBase) error {
	tq.typeQuery[funcName] = &TypeQuery{
		typeQuery: WRITE,
		cached:    time.Now(),
	}
	t, err := getQueryType(funcName, db)
	if err != nil {
		return err
	}
	tq.typeQuery[funcName].typeQuery = t
	return nil
}

func (tq *TypeQueryPool) isCached(funcName string) bool {
	if tq.typeQuery[funcName] == nil {
		return false
	}
	return tq.typeQuery[funcName].cached.Add(cacheTimeLife).After(time.Now())
}

func getQueryFunction(query string) (string, error) {
	start := strings.Index(query, "FROM ")
	if start == -1 {
		start = strings.Index(query, "SELECT ")
		if start == -1 {
			return "", indexError
		}
		start += 7
	} else {
		start += 5
	}
	end := strings.Index(query[start:], "(")
	if end == -1 {
		return "", indexError
	}
	return query[start : start+end], nil
}

func getQueryType(fullFuncName string, db DataBase) (string, error) {
	if db == nil {
		return "", dataBaseError
	}

	names := strings.Split(fullFuncName, ".")
	if len(names) != 2 {
		names = append([]string{"public"}, names...)
	}
	funcType := ""
	err := db.QueryRow(
		"SELECT provolatile FROM pg_proc JOIN pg_namespace ON pg_proc.pronamespace = pg_namespace.oid WHERE proname=$1 AND pg_namespace.nspname = $2",
		names[1],
		names[0],
	).Scan(&funcType)
	if err != nil {
		return "", err
	}
	if (funcType == "s") || (funcType == "i") {
		return READ, nil
	}
	return WRITE, nil
}
