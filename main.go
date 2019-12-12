package gol

const (
	DatabaseTypePostgresql = "postgres"
	DatabaseTypeMysql      = "mysql"
)

func Open(databaseType string, host string, port string, user string, pass string, database string, optionMap map[string]string) (*DB, error) {
	var err error

	db := &DB{
		DB:       nil,
		TX:       nil,
		modeLog:  false,
		modeTest: false,
	}

	err = db.Init(databaseType, host, port, user, pass, database, optionMap)
	if err != nil {
		return nil, err
	}

	return db, nil
}
