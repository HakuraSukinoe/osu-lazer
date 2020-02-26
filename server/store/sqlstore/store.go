package sqlstore

import (
	"github.com/deissh/osu-lazer/server/store"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SqlStore interface {
	DriverName() string
	GetCurrentSchemaVersion() string
	GetMaster() *sqlx.DB
	TotalMasterDbConnections() int

	User() store.UserStore
}
