package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// No Surfs protexts Same site issues in Post requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandller := nosurf.New(next)
	csrfHandller.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.IsProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandller

}

// Loads session request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
