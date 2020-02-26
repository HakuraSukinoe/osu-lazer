package model

const (
	SQL_SETTINGS_DEFAULT_DATA_SOURCE = "postgres://postgres:postgres@/osuserver?sslmode=disable"
	DATABASE_NAME                    = "osuserver"
	DATABASE_DRIVER_POSTGRES         = "postgres"
)

type SqlSettings struct {
	DriverName                  *string `restricted:"true"`
	DataSource                  *string `restricted:"true"`
	DatabaseName                *string `restricted:"true"`
	MaxIdleConns                *int    `restricted:"true"`
	ConnMaxLifetimeMilliseconds *int    `restricted:"true"`
	MaxOpenConns                *int    `restricted:"true"`
	Trace                       *bool   `restricted:"true"`
	AtRestEncryptKey            *string `restricted:"true"`
	QueryTimeout                *int    `restricted:"true"`
}

func (s *SqlSettings) SetDefaults(isUpdate bool) {
	if s.DriverName == nil {
		s.DriverName = NewString(DATABASE_DRIVER_POSTGRES)
	}

	if s.DataSource == nil {
		s.DataSource = NewString(SQL_SETTINGS_DEFAULT_DATA_SOURCE)
	}

	if s.DatabaseName == nil {
		s.DatabaseName = NewString(DATABASE_NAME)
	}

	if isUpdate {
		// When updating an existing configuration, ensure an encryption key has been specified.
		if s.AtRestEncryptKey == nil || len(*s.AtRestEncryptKey) == 0 {
			s.AtRestEncryptKey = NewString(NewRandomString(32))
		}
	} else {
		// When generating a blank configuration, leave this key empty to be generated on server start.
		s.AtRestEncryptKey = NewString("")
	}

	if s.MaxIdleConns == nil {
		s.MaxIdleConns = NewInt(20)
	}

	if s.MaxOpenConns == nil {
		s.MaxOpenConns = NewInt(300)
	}

	if s.ConnMaxLifetimeMilliseconds == nil {
		s.ConnMaxLifetimeMilliseconds = NewInt(3600000)
	}

	if s.Trace == nil {
		s.Trace = NewBool(false)
	}

	if s.QueryTimeout == nil {
		s.QueryTimeout = NewInt(30)
	}
}
