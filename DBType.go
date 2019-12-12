package gol

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	// Not Use ssl, Not check the certificate.
	DatabaseSslModeDisable = "disable"
	// Use ssl, Not check the certificate.
	DatabaseSslModeRequire = "require"
	// Use ssl, Check the certificate
	DatabaseSslModeVerifyCa = "verify-ca"
	// Use ssl, Check the certificate and confirm that it is on the server.
	DatabaseSslModeVerifyFull = "verify-full"
)

type DB struct {
	DB               *sql.DB
	TX               *sql.Tx
	modeDatabaseType string
	modeLog          bool
	modeResultKey    int
	modeTest         bool
}

func (rec *DB) Init(databaseType string, host string, port string, user string, pass string, database string, optionMap map[string]string) error {
	var err error
	var db *sql.DB

	switch databaseType {
	case DatabaseTypeMysql:
		// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		source := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			user,
			pass,
			host,
			port,
			database,
		)

		db, err = sql.Open(DatabaseTypeMysql, source)
		if err != nil {
			return err
		}

		rec.modeDatabaseType = DatabaseTypeMysql
	case DatabaseTypePostgresql:

		sslMode := DatabaseSslModeDisable
		if v, ok := optionMap["sslMode"]; ok {
			switch v {
			case DatabaseSslModeDisable:
				sslMode = DatabaseSslModeDisable
			case DatabaseSslModeRequire:
				sslMode = DatabaseSslModeRequire
			case DatabaseSslModeVerifyCa:
				sslMode = DatabaseSslModeVerifyCa
			case DatabaseSslModeVerifyFull:
				sslMode = DatabaseSslModeVerifyFull
			}
		}

		source := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host,
			port,
			user,
			pass,
			database,
			sslMode,
		)

		db, err = sql.Open(DatabaseTypePostgresql, source)
		if err != nil {
			return err
		}

		rec.modeDatabaseType = DatabaseTypePostgresql
	default:
		return errors.New("unknown databaseType")
	}

	rec.DB = db
	rec.SetModeResultKey()
	rec.SetModeLog(false)

	return nil
}

func (rec *DB) SetModeLog(mode bool) {
	rec.modeLog = mode
}

func (rec *DB) SetModeResultKey() {
	rec.modeResultKey = resultKeyModeNone
}

func (rec *DB) SetModeResultKeyCamelCase() {
	rec.modeResultKey = resultKeyModeCamelCase
}

func (rec *DB) SetModeResultKeySnakeCase() {
	rec.modeResultKey = resultKeyModeSnakeCase
}

func (rec *DB) Query() *QueryType {
	queryData := &QueryType{}

	queryData.Init(rec.DB, rec.TX, rec.modeDatabaseType)
	queryData.SetModeLog(rec.modeLog)

	switch rec.modeResultKey {
	case resultKeyModeCamelCase:
		queryData.SetModeResultKeyCamelCase()
	case resultKeyModeSnakeCase:
		queryData.SetModeResultKeySnakeCase()
	default:
		queryData.SetModeResultKey()
	}

	return queryData
}

func (rec *DB) Exec(query string, valueList ...interface{}) (sql.Result, error) {
	queryData := rec.Query()

	result, err := queryData.Exec(query, valueList...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (rec *DB) ExecQuery(dest interface{}, query string, valueList ...interface{}) error {
	queryData := rec.Query()

	err := queryData.ExecQuery(dest, query, valueList...)
	if err != nil {
		return err
	}

	return nil
}

func (rec *DB) Close() error {
	if rec.modeTest {
		return nil
	}

	if rec.DB == nil {
		return nil
	}

	err := rec.DB.Close()
	if err != nil {
		return err
	}

	rec.DB = nil
	rec.TX = nil

	return nil
}

func (rec *DB) Begin() (*DB, error) {
	if rec.modeTest {
		return rec, nil
	}

	if rec.TX != nil {
		return nil, errors.New("exist transaction")
	}

	tx, err := rec.DB.Begin()
	if err != nil {
		return nil, err
	}

	queryData := *rec
	queryData.TX = tx

	return &queryData, nil
}

func (rec *DB) Commit() error {
	var err error

	if rec.modeTest {
		return nil
	}

	if rec.TX == nil {
		return nil
	}

	err = rec.TX.Commit()
	if err != nil {
		return err
	}

	rec.TX = nil

	return nil
}

func (rec *DB) Rollback() error {
	var err error

	if rec.modeTest {
		return nil
	}

	if rec.TX == nil {
		return nil
	}

	err = rec.TX.Rollback()
	if err != nil {
		return err
	}

	rec.TX = nil

	return nil
}

func (rec *DB) TestStart() {
	if rec.DB == nil {
		return
	}

	tx, _ := rec.DB.Begin()
	rec.modeTest = true
	rec.TX = tx
}

func (rec *DB) TestEnd() {
	var err error

	err = rec.TX.Rollback()
	if err != nil {
		panic(err)
	}

	rec.TX = nil

	_ = rec.DB.Close()
	rec.DB = nil
}
