package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pradeepj4u/bookings/internal/config"
	"github.com/pradeepj4u/bookings/internal/handler"
)

func routes(app *config.AppConfig) http.Handler {
	// mux := pat.New()
	// mux.Get("/", http.HandlerFunc(handler.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handler.Repo.About))

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handler.Repo.Home)

	mux.Get("/about", handler.Repo.About)

	mux.Get("/standard-room", handler.Repo.StandardRoom)
	mux.Get("/king-suit", handler.Repo.KingSuit)

	mux.Get("/make-reservations", handler.Repo.MakeReservations)
	mux.Post("/make-reservations", handler.Repo.PostMakeReservations)
	mux.Get("/reservation-summary", handler.Repo.ReservationSummary)

	mux.Get("/check-availability", handler.Repo.CheckAvailability)
	mux.Post("/check-availability", handler.Repo.PostCheckAvailability)
	mux.Post("/check-availability-json", handler.Repo.CheckAvailabilityJson)

	mux.Get("/book-room", handler.Repo.BookRoom)

	mux.Get("/user/login", handler.Repo.ShowLogin)
	mux.Post("/user/login", handler.Repo.PostLogin)

	mux.Get("/user/logout", handler.Repo.Logout)

	mux.Get("/choose-room/{roomid}", handler.Repo.ChooseRoom)

	mux.Get("/contact", handler.Repo.ContactUs)

	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth)
		mux.Get("/dashboard", handler.Repo.AdminDashboard)

		mux.Get("/reservations-new", handler.Repo.AdminNewReservations)
		mux.Get("/reservations-all", handler.Repo.AdminAllReservations)
		
		mux.Get("/process-reservation/{src}/{id}/do", handler.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}/do", handler.Repo.AdminDeleteReservation)

		mux.Get("/reservation/{src}/{id}/show", handler.Repo.AdminShowReservation)
		mux.Post("/reservation/{src}/{id}", handler.Repo.AdminPostShowReservation)

		mux.Get("/reservation-calendar", handler.Repo.AdminReservationsCalendar)
		mux.Post("/reservation-calendar", handler.Repo.AdminPostReservationsCalendar)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
