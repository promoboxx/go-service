package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"sync"

	"github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"
	"github.com/promoboxx/go-discovery/src/discovery"
	"github.com/promoboxx/go-metric-client/metrics"
)

func init() {
	sql.Register("instrumented-postgres", instrumentedsql.WrapDriver(&pq.Driver{},
		instrumentedsql.WithTracer(opentracing.NewTracer(false)),
		instrumentedsql.WithOmitArgs()))
}

type SQLDBConnector interface {
	GetConnection() (*sql.DB, error)
}

type sqlDBConnector struct {
	maxOpenConns  int
	maxIdleConns  int
	dbName        string
	driver        string
	dbUser        string
	dbPass        string
	sslMode       string
	finder        discovery.Finder
	metricsClient metrics.Client
}

// keeps track of the db pool for each connection string
// if the db moves and new connection string comes in it will create a new pool, otherwise we keep the existing one
var connMap map[string]*sql.DB
var mapLock sync.RWMutex

func init() {
	mapLock.Lock()
	defer mapLock.Unlock()
	connMap = make(map[string]*sql.DB)
}

func NewSQLDBConnector(maxOpenConns, maxIdleConns int, dbName string, driver string, sslMode string, dbUser string, dbPass string, finder discovery.Finder, metricsClient metrics.Client) SQLDBConnector {
	return &sqlDBConnector{
		maxOpenConns:  maxOpenConns,
		maxIdleConns:  maxIdleConns,
		dbName:        dbName,
		driver:        driver,
		dbUser:        dbUser,
		dbPass:        dbPass,
		sslMode:       sslMode,
		finder:        finder,
		metricsClient: metricsClient,
	}
}

func (c *sqlDBConnector) GetConnection() (*sql.DB, error) {
	connString, err := c.getConnectionString()
	if connString == "" || err != nil {
		return nil, fmt.Errorf("Could not find connection string for service: %v- %s", c.dbName, err)
	}

	db, err := c.getDbFromMap(c.driver, connString)
	if err != nil {
		return nil, fmt.Errorf("Could not get db connection: %v", err)
	}

	if c.metricsClient != nil {
		c.metricsClient.InternalCustom("database.connections", "get-connection", "connection-pool-max-conns", map[string]string{"db-name": c.dbName}, int64(db.Stats().MaxOpenConnections))
		c.metricsClient.InternalCustom("database.connections", "get-connection", "connection-pool-max-idle-conns", map[string]string{"db-name": c.dbName}, int64(db.Stats().Idle))
		c.metricsClient.InternalCustom("database.connections", "get-connection", "connection-pool-open-conns", map[string]string{"db-name": c.dbName}, int64(db.Stats().OpenConnections))
		c.metricsClient.InternalCustom("database.connections", "get-connection", "connection-pool-in-use-conns", map[string]string{"db-name": c.dbName}, int64(db.Stats().InUse))
	}

	return db, nil
}

// GetDbFromMap is a thread safe way to get a connection from the map or create a new one
func (c *sqlDBConnector) getDbFromMap(driver string, conn string) (result *sql.DB, err error) {
	key := fmt.Sprintf("%s|%s", c.driver, conn)
	var ok bool
	mapLock.RLock()
	result, ok = connMap[key]
	mapLock.RUnlock()
	if ok {
		return
	}
	// need to create a connection
	mapLock.Lock()
	defer mapLock.Unlock()
	result, ok = connMap[key] // check to make sure something didn't just create it
	if ok {
		return
	}

	result, err = sql.Open(driver, conn)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to DB: %v", err)
	}
	result.SetMaxOpenConns(c.maxOpenConns)
	result.SetMaxIdleConns(c.maxIdleConns)
	connMap[key] = result
	return
}

func (c *sqlDBConnector) getConnectionString() (string, error) {
	dbAddr, dbPort, err := c.finder.FindHostPort(fmt.Sprintf("%s-db", c.dbName))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.dbUser, url.QueryEscape(c.dbPass), dbAddr, dbPort, c.dbName, c.sslMode), nil
}
