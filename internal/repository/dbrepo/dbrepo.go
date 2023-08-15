package dbrepo

import (
	"database/sql"

	"github.com/pradeepj4u/bookings/internal/config"
	"github.com/pradeepj4u/bookings/internal/repository"
)

// class for creating Postgress connections
type postgresDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// constructor for posgress connection class
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepo{
		App: a,
		DB:  conn,
	}
}
