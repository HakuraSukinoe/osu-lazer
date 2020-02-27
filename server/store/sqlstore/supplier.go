package sqlstore

import (
	"context"
	"github.com/deissh/osu-lazer/server/model"
	"github.com/deissh/osu-lazer/server/store"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

const (
	DB_PING_ATTEMPTS     = 18
	DB_PING_TIMEOUT_SECS = 10
)

const (
	EXIT_GENERIC_FAILURE = 1
	EXIT_DB_OPEN         = 101
	EXIT_PING            = 102
)

type SqlSupplierStores struct {
	user   store.UserStore
	friend store.FriendStore
}

type SqlSupplier struct {
	rrCounter      int64
	srCounter      int64
	master         *sqlx.DB
	stores         SqlSupplierStores
	settings       *model.SqlSettings
	lockedToMaster bool
}

func NewSqlSupplier(settings *model.SqlSettings) *SqlSupplier {
	supplier := &SqlSupplier{
		rrCounter: 0,
		srCounter: 0,
		settings:  settings,
	}
	supplier.initConnection()

	supplier.stores.user = NewSqlUserStore(supplier)

	return supplier
}

func (ss *SqlSupplier) initConnection() {
	db, err := sqlx.Connect(ss.settings.DriverName, ss.settings.DataSource)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to open SQL connection to err.")
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
				log.Fatal().
					Err(err).
					Msg("Failed to ping DB, server will exit.")
				time.Sleep(time.Second)
				os.Exit(EXIT_PING)
			} else {
				log.Error().
					Err(err).
					Int("timeout", DB_PING_TIMEOUT_SECS).
					Msg("Failed to ping DB")
				time.Sleep(DB_PING_TIMEOUT_SECS * time.Second)
			}
		}
	}

	db.SetMaxIdleConns(ss.settings.MaxIdleConns)
	db.SetMaxOpenConns(ss.settings.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(ss.settings.ConnMaxLifetimeMilliseconds) * time.Millisecond)

	ss.master = db
}

func (ss *SqlSupplier) DriverName() string {
	return ss.settings.DriverName
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

func (ss *SqlSupplier) User() store.UserStore {
	return ss.stores.user
}

func (ss *SqlSupplier) Friend() store.FriendStore {
	return ss.stores.friend
}
