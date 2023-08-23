package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/pradeepj4u/bookings/internal/helpers"
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

// Handels Authrntication
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuntenticated(r) {
			session.Put(r.Context(), "CriticalEdit", "Login First")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
