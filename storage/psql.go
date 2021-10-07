package storage

import (
	"database/sql"
)

//PostgreSQL struct for postgresql database operations
type PostgreSQL struct {
	db *sql.DB
}

//Init - open database
func (psql *PostgreSQL) Init(connectionString string) error {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	psql.db = db

	return nil
}
