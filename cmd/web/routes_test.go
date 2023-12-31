package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pradeepj4u/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do noting pass
	default:
		t.Errorf("Incorrect return type of routes %T", v)
	}
}
