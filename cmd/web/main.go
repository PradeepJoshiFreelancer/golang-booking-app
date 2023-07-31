package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pradeepj4u/bookings/pkg/config"
	"github.com/pradeepj4u/bookings/pkg/handler"
	"github.com/pradeepj4u/bookings/pkg/render"
)

var portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {

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

	repo := handler.NewRpository(&app)
	handler.NewHandller(repo)

	render.CreateNewAppConfig(&app)

	svr := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	log.Println("Starting at port:", portNumber)
	erro := svr.ListenAndServe()
	log.Println(erro)
}