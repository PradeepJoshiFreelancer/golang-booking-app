package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pradeepj4u/bookings/cmd/models"
	"github.com/pradeepj4u/bookings/driver"
	"github.com/pradeepj4u/bookings/internal/config"
	"github.com/pradeepj4u/bookings/internal/handler"
	"github.com/pradeepj4u/bookings/internal/helpers"
	"github.com/pradeepj4u/bookings/internal/render"
)

var portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)
	//listen to the email
	listenToMail()

	log.Println("Starting at port:", portNumber)
	svr := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	erro := svr.ListenAndServe()
	log.Fatal(erro)
}

func run() (*driver.DB, error) {

	// What value we are going to save in the context.
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(models.Reservation{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	//read the flag from Command line
	inProduction := flag.Bool("production", true, "Application is in production.")
	useCache := flag.Bool("cache", true, "Cache templets")
	dbName := flag.String("dbname", "", "Database Name")
	dbHost := flag.String("dbhost", "localhost", "Database Host Name")
	dbUser := flag.String("dbuser", "", "Database User")
	dbPassword := flag.String("dbpass", "", "Database Password")
	dbPort := flag.String("dbport", "5432", "Database Port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings(disable, prefer, required)")

	flag.Parse()

	if *dbName == "" || *dbUser == ""{
		fmt.Println("Missing required flags dbName or dbUser")
		os.Exit(1)
	}

	emailData := make(chan models.EmailData, 10)
	app.MailChan = emailData

	//changes IsProduction
	app.IsProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	//making a database connection
	log.Println("Connecting to the database..")
	//"host=localhost port=5432 dbname=bookings user=postgres password=India@100"
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPassword, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	log.Println("Connected to the database!")

	tc, err := render.CreateChacheMap()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	app.TempletCache = tc

	repo := handler.NewRpository(&app, db)
	handler.NewHandller(repo)

	render.NewRenderer(&app)

	helpers.NewHelper(&app)

	return db, nil
}
