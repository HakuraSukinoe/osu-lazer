package sqlstore

import (
	"context"
	"github.com/deissh/osu-lazer/server/einterfaces"
	"github.com/deissh/osu-lazer/server/mlog"
	"github.com/deissh/osu-lazer/server/model"
	"github.com/deissh/osu-lazer/server/store"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"os"
	"strings"
	"time"
)

const (
	DB_PING_ATTEMPTS     = 18
	DB_PING_TIMEOUT_SECS = 10
)

const (
	EXIT_GENERIC_FAILURE             = 1
	EXIT_CREATE_TABLE                = 100
	EXIT_DB_OPEN                     = 101
	EXIT_PING                        = 102
	EXIT_NO_DRIVER                   = 103
	EXIT_TABLE_EXISTS                = 104
	EXIT_TABLE_EXISTS_MYSQL          = 105
	EXIT_COLUMN_EXISTS               = 106
	EXIT_DOES_COLUMN_EXISTS_POSTGRES = 107
	EXIT_DOES_COLUMN_EXISTS_MYSQL    = 108
	EXIT_DOES_COLUMN_EXISTS_MISSING  = 109
	EXIT_CREATE_COLUMN_POSTGRES      = 110
	EXIT_CREATE_COLUMN_MYSQL         = 111
	EXIT_CREATE_COLUMN_MISSING       = 112
	EXIT_REMOVE_COLUMN               = 113
	EXIT_RENAME_COLUMN               = 114
	EXIT_MAX_COLUMN                  = 115
	EXIT_ALTER_COLUMN                = 116
	EXIT_CREATE_INDEX_POSTGRES       = 117
	EXIT_CREATE_INDEX_MYSQL          = 118
	EXIT_CREATE_INDEX_FULL_MYSQL     = 119
	EXIT_CREATE_INDEX_MISSING        = 120
	EXIT_REMOVE_INDEX_POSTGRES       = 121
	EXIT_REMOVE_INDEX_MYSQL          = 122
	EXIT_REMOVE_INDEX_MISSING        = 123
	EXIT_REMOVE_TABLE                = 134
	EXIT_CREATE_INDEX_SQLITE         = 135
	EXIT_REMOVE_INDEX_SQLITE         = 136
	EXIT_TABLE_EXISTS_SQLITE         = 137
	EXIT_DOES_COLUMN_EXISTS_SQLITE   = 138
)

type SqlSupplierStores struct {
	user store.UserStore
}

type SqlSupplier struct {
	// rrCounter and srCounter should be kept first.
	// See https://github.com/mattermost/mattermost-server/v5/pull/7281
	rrCounter      int64
	srCounter      int64
	master         *sqlx.DB
	stores         SqlSupplierStores
	settings       *model.SqlSettings
	lockedToMaster bool
}

func NewSqlSupplier(settings model.SqlSettings, metrics einterfaces.MetricsInterface) *SqlSupplier {
	supplier := &SqlSupplier{
		rrCounter: 0,
		srCounter: 0,
		settings:  &settings,
	}
	supplier.initConnection()

	supplier.stores.user = NewSqlUserStore(supplier, metrics)

	return supplier
}

func (ss *SqlSupplier) initConnection() {
	db, err := sqlx.Connect(*ss.settings.DriverName, *ss.settings.DataSource)
	if err != nil {
		mlog.Critical("Failed to open SQL connection to err.", mlog.Err(err))
		time.Sleep(time.Second)
		os.Exit(EXIT_DB_OPEN)
	}

	for i := 0; i < DB_PING_ATTEMPTS; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), DB_PING_TIMEOUT_SECS*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err == nil {
			break
		} else {
			if i == DB_PING_ATTEMPTS-1 {
				mlog.Critical("Failed to ping DB, server will exit.", mlog.Err(err))
				time.Sleep(time.Second)
				os.Exit(EXIT_PING)
			} else {
				mlog.Error("Failed to ping DB", mlog.Err(err), mlog.Int("retrying in seconds", DB_PING_TIMEOUT_SECS))
				time.Sleep(DB_PING_TIMEOUT_SECS * time.Second)
			}
		}
	}

	db.SetMaxIdleConns(*ss.settings.MaxIdleConns)
	db.SetMaxOpenConns(*ss.settings.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(*ss.settings.ConnMaxLifetimeMilliseconds) * time.Millisecond)

	ss.master = db
}

func (ss *SqlSupplier) DriverName() string {
	return *ss.settings.DriverName
}

func (ss *SqlSupplier) GetCurrentSchemaVersion() string {
	var version string
	_ = ss.GetMaster().Get(&version, "SELECT version FROM schema_migrations")

	return version
}

func (ss *SqlSupplier) GetMaster() *sqlx.DB {
	return ss.master
}

func (ss *SqlSupplier) TotalMasterDbConnections() int {
	return ss.GetMaster().Stats().OpenConnections
}

func IsUniqueConstraintError(err error, indexName []string) bool {
	unique := false
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		unique = true
	}

	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		unique = true
	}

	field := false
	for _, contain := range indexName {
		if strings.Contains(err.Error(), contain) {
			field = true
			break
		}
	}

	return unique && field
}

func (ss *SqlSupplier) Close() {
	_ = ss.master.Close()
}

func (ss *SqlSupplier) LockToMaster() {
	ss.lockedToMaster = true
}

func (ss *SqlSupplier) UnlockFromMaster() {
	ss.lockedToMaster = false
}

func (ss *SqlSupplier) User() store.UserStore {
	return ss.stores.user
}

type JSONSerializable interface {
	ToJson() string
}

func convertMySQLFullTextColumnsToPostgres(columnNames string) string {
	columns := strings.Split(columnNames, ", ")
	concatenatedColumnNames := ""
	for i, c := range columns {
		concatenatedColumnNames += c
		if i < len(columns)-1 {
			concatenatedColumnNames += " || ' ' || "
		}
	}

	return concatenatedColumnNames
}
