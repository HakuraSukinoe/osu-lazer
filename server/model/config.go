package model

import "github.com/gookit/config/v2"

const (
	DEFAULT_DSN              = "postgres://postgres:postgres@/osuserver?sslmode=disable"
	DATABASE_DRIVER_POSTGRES = "postgres"

	MAX_IDLE_CONNS    = 20
	MAX_OPEN_CONNS    = 300
	CONN_MAX_LIFETIME = 3600000
	QUERY_TIMEOUT     = 30
)

type SqlSettings struct {
	DriverName                  string
	DataSource                  string
	MaxIdleConns                int
	ConnMaxLifetimeMilliseconds int
	MaxOpenConns                int
	Trace                       bool
	AtRestEncryptKey            string
	QueryTimeout                int
}

func NewSqlSettings() *SqlSettings {
	s := SqlSettings{}

	s.DriverName = config.String("server.database.driver", DATABASE_DRIVER_POSTGRES)
	s.DataSource = config.String("server.database.dsn", DEFAULT_DSN)
	s.MaxIdleConns = config.Int("server.database.max_idle_conns", MAX_IDLE_CONNS)
	s.MaxOpenConns = config.Int("server.database.max_open_conns", MAX_OPEN_CONNS)
	s.ConnMaxLifetimeMilliseconds = config.Int("server.database.conn_max_lifetime", CONN_MAX_LIFETIME)
	s.Trace = config.Bool("server.database.trace", false)
	s.QueryTimeout = config.Int("server.database.query_timeout", QUERY_TIMEOUT)

	return &s
}
