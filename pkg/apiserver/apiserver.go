package apiserver

import (
	"github.com/sirupsen/logrus"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/tools"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/tools/cluster"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/tools/postgres"
	"time"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/config"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store/sqlstore"
)

func Start(isClusterDb bool) error {

	cfg := config.GetInstance()

	db, err := connectDB(isClusterDb)

	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	server := NewServer(store)
	logrus.Infof("%v API Server 'Login-service' is started in addr:[%s]", time.Now(), cfg.BindAddr)

	return server.Start(cfg.BindAddr)
}

func connectDB(isClusterDb bool) (tools.DataBase, error) {
	if isClusterDb {
		return cluster.Start()

	} else {
		db, err := postgres.OpenDataBase("postgres", config.GetInstance().DatabaseURL, &postgres.DataBaseParams{EnableLogging: true})
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			return nil, err
		}

		return db, nil
	}
}
