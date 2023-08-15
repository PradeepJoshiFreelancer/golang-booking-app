package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

var maxOpenConn = 10
var maxIdelDBConn = 5

const maxConnLifetime = 5 * time.Minute

// ConnectSQL makes connection to a datbase
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenConn)
	d.SetMaxIdleConns(maxIdelDBConn)
	d.SetConnMaxLifetime(maxConnLifetime)
	dbConn.SQL = d
	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// textDB sends a [ing to the database to check if db is active]
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// creates a new database
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
