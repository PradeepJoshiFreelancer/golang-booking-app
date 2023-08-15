package handler

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/pradeepj4u/bookings/cmd/models"
	"github.com/pradeepj4u/bookings/internal/config"
	"github.com/pradeepj4u/bookings/internal/render"
)

var session *scs.SessionManager
var app config.AppConfig
var pathToTemplet = "./../../templet/"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	// What value we are going to save in the context.
	gob.Register(models.Reservation{})

	//changes IsProduction
	app.IsProduction = false

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	// fmt.Println("Hello World")
	tc, err := render.CreateChacheMap()
	if err != nil {
		log.Fatal(err)
	}
	app.TempletCache = tc
	app.UseCache = true

	repo := NewRpository(&app)
	NewHandller(repo)

	render.NewRenderer(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)

	mux.Get("/about", Repo.About)

	mux.Get("/standard-room", Repo.StandardRoom)
	mux.Get("/king-suit", Repo.KingSuit)

	mux.Get("/make-reservations", Repo.MakeReservations)
	mux.Post("/make-reservations", Repo.PostMakeReservations)

	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/check-availability", Repo.CheckAvailability)
	mux.Post("/check-availability", Repo.PostCheckAvailability)
	mux.Post("/check-availability-json", Repo.CheckAvailabilityJson)

	mux.Get("/contact", Repo.ContactUs)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

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

func CreateTestChacheMap() (map[string]*template.Template, error) {
	chacheMap := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplet))
	if err != nil {
		return chacheMap, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page)

		ts, err := template.New(fileName).Funcs(functions).ParseFiles(page)

		if err != nil {

			return chacheMap, err
		}
		layout, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplet))
		if err != nil {
			return chacheMap, err
		}
		if len(layout) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplet))
			if err != nil {

				return chacheMap, err
			}
		}
		chacheMap[fileName] = ts
	}

	return chacheMap, nil

}
