package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// create Application Templet that is available cross program
type AppConfig struct {
	UseCache     bool
	TempletCache map[string]*template.Template
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	IsProduction bool
	Session      *scs.SessionManager
}
