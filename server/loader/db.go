package loader

import (
	"database/sql"
)

func ConnectDB(dbDriver, dbSource string) (*sql.DB, error) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return &sql.DB{}, err
	}

	conn.SetMaxIdleConns(2)
	conn.SetMaxOpenConns(2)

	return conn, nil
}
