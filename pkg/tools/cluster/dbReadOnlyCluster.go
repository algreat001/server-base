package cluster

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/config"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/tools"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/tools/postgres"
)

type DbCluster struct {
	mainDB             tools.DataBase
	readOnlyDBs        []tools.DataBase
	currentServerIndex int
	typeQueryCache     *tools.TypeQueryPool
}

type DbClusterStats struct {
	Pid             int
	UseSysID        int
	UserName        string
	ApplicationName string
	ClientAddr      string
	ClientHostname  string
	ClientPort      int
	BackendStart    time.Time
	BackendXmin     *int
	State           string
	SentLSN         string
	WriteLSN        string
	FlushLSN        string
	ReplayLSN       string
	WriteLag        string
	FlushLag        string
	ReplayLag       string
	SyncPriority    int
	SyncState       string
	ReplyTime       time.Time
}

func Start() (tools.DataBase, error) {
	typeQueryCache := tools.NewTypeQueryPool()

	mainDb, err := postgres.OpenDataBaseWithServerName("postgres", "main", config.GetInstance().Database.MainServer, &postgres.DataBaseParams{EnableLogging: true})
	if err != nil {
		return nil, err
	}

	readOnlyDbs := make([]tools.DataBase, 0)
	for inx, server := range config.GetInstance().Database.ReadServers {
		roDb, err := postgres.OpenDataBaseWithServerName("postgres", fmt.Sprintf("replica-%d", inx), server, &postgres.DataBaseParams{EnableLogging: true})
		if err != nil {
			return nil, err
		}
		readOnlyDbs = append(readOnlyDbs, roDb)
	}

	dbCluster := &DbCluster{
		mainDB:             mainDb,
		readOnlyDBs:        readOnlyDbs,
		currentServerIndex: 0,
		typeQueryCache:     typeQueryCache,
	}

	if err := dbCluster.Ping(); err != nil {
		return nil, err
	}

	return dbCluster, nil
}

func (d *DbCluster) Close() error {
	err := d.mainDB.Close()
	if err != nil {
		return err
	}
	for _, db := range d.readOnlyDBs {
		err := db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DbCluster) Exec(query string, args ...any) (sql.Result, error) {
	return d.getDB(query).Exec(query, args...)
}

func (d *DbCluster) Query(query string, args ...any) (*sql.Rows, error) {
	return d.getDB(query).Query(query, args...)
}

func (d *DbCluster) QueryRow(query string, args ...any) *sql.Row {
	return d.getDB(query).QueryRow(query, args...)
}

func (d *DbCluster) Ping() error {
	err := d.mainDB.Ping()
	if err != nil {
		return err
	}
	for _, db := range d.readOnlyDBs {
		err := db.Ping()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DbCluster) GetStats() string {
	return "unimplemented"
}

// getDB returns the database connection to use for the query.
func (d *DbCluster) getDB(query string) tools.DataBase {
	if d.typeQueryCache.GetTypeQuery(query, d.mainDB) == tools.WRITE {
		return d.mainDB
	}
	d.currentServerIndex++ // round-robin
	if d.currentServerIndex >= len(d.readOnlyDBs) {
		d.currentServerIndex = -1
	}
	if d.currentServerIndex == -1 {
		return d.mainDB
	}
	return d.readOnlyDBs[d.currentServerIndex]
}
