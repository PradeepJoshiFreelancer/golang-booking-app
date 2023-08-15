package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pradeepj4u/bookings/cmd/models"
	"github.com/pradeepj4u/bookings/driver"
	"github.com/pradeepj4u/bookings/internal/config"
	"github.com/pradeepj4u/bookings/internal/forms"
	"github.com/pradeepj4u/bookings/internal/helpers"
	"github.com/pradeepj4u/bookings/internal/render"
	"github.com/pradeepj4u/bookings/internal/repository"
	"github.com/pradeepj4u/bookings/internal/repository/dbrepo"
)

// Stores the Class address
var Repo *Repository

// New Class Repository
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Creates a new instance of class Repository
func NewRpository(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// sets the Repo Variable for new handllers
func NewHandller(r *Repository) {
	Repo = r
}

// Renders the Home Page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	render.ParseTemplet(w, r, "home.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}

// Renders the About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := map[string]string{}

	render.ParseTemplet(w, r, "about.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}

// Renders the Room page
func (m *Repository) StandardRoom(w http.ResponseWriter, r *http.Request) {
	render.ParseTemplet(w, r, "standard-room.page.tmpl", &models.TempletData{})
}

// Renders the Check availability post
func (m *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	render.ParseTemplet(w, r, "check-availability.page.tmpl", &models.TempletData{})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomId    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Returns a JSON response
func (m *Repository) CheckAvailabilityJson(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	// Parse 2006-02-01 -> 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomId, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailibilityForDateRangeByRoomId(roomId, startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "This is the response Message",
		RoomId:    r.Form.Get("room_id"),
		StartDate: sd,
		EndDate:   ed,
	}
	out, err := json.MarshalIndent(resp, "", "      ")
	if err != nil {
		m.App.ErrorLog.Println("Can not parse the response Json")
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Take Query parms create session and routes to make reservations
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	startDate := r.URL.Query().Get("s")
	endDate := r.URL.Query().Get("e")

	var res models.Reservation

	room, err := m.DB.GetRoomById(roomId)
	if err != nil {
		helpers.ServerError(w, err)
	}
	res.Room = room
	res.RoomId = roomId

	layout := "2006-01-02"
	sd, err := time.Parse(layout, startDate)
	if err != nil {
		helpers.ServerError(w, err)
	}
	ed, err := time.Parse(layout, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.StartDate = sd
	res.EndDate = ed

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)
}

// Check availability post
func (m *Repository) PostCheckAvailability(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	// Parse 2006-02-01 -> 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailibilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "CriticalEdit", "No availibility for the dates selected.")
		http.Redirect(w, r, "/check-availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.ParseTemplet(w, r, "choose-room.page.tmpl", &models.TempletData{
		Data: data,
	})

}

// Renders the Room page
func (m *Repository) KingSuit(w http.ResponseWriter, r *http.Request) {
	render.ParseTemplet(w, r, "king-suit.page.tmpl", &models.TempletData{})
}

// Renders the Contact Us page
func (m *Repository) ContactUs(w http.ResponseWriter, r *http.Request) {
	render.ParseTemplet(w, r, "contact.page.tmpl", &models.TempletData{})
}

// Renders the Check availability post
func (m *Repository) MakeReservations(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot find reservations data in the sessions"))
		return
	}
	data := make(map[string]interface{})

	data["reservation"] = res

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	m.App.Session.Put(r.Context(), "reservation", res)

	render.ParseTemplet(w, r, "make-reservations.page.tmpl", &models.TempletData{
		Form:      forms.NewForm(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// Takes input from make reservation page and routes to summary pages
func (m *Repository) PostMakeReservations(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot find reservations data in the sessions"))
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	newForm := forms.NewForm(r.PostForm)

	newForm.Required("first_name", "last_name", "email", "phone")
	newForm.MinLength("first_name", 3, r)
	newForm.IsEmailValid("email")

	if !newForm.IsFormValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.ParseTemplet(w, r, "make-reservations.page.tmpl", &models.TempletData{
			Form: newForm,
			Data: data,
		})
		return
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)
	newReservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}
	m.DB.InsertRoomRestriction(restriction)

	//Send Notifications
	content := fmt.Sprintf(`
	<strong>Booking Confirmed</strong><br>
	Dear %s,<br>
	Your booking in hotal is confirmed from %s to %s.

	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.EmailData{
		To:      reservation.Email,
		From:    "hotel@hotel.com",
		Subject: "Booking Confirmation",
		Content: content,
		Templet: "basic.html",
	}

	m.App.MailChan <- msg

	content = fmt.Sprintf(`
	<strong>New Booking</strong><br>
	Dear Property owner,<br>
	Your booking in hotal is confirmed from %s to %s.

	`, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg = models.EmailData{
		To:      reservation.Email,
		From:    "hotel@hotel.com",
		Subject: "New Booking",
		Content: content,
		Templet: "basic.html",
	}

	m.App.MailChan <- msg

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Renders the Reservation Summary Page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservationPageData, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "CriticalEdit", "Can't get reservation in Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		m.App.ErrorLog.Println("Unable to get data from Session")
		return
	}
	stringMap := make((map[string]string))
	stringMap["start_date"] = reservationPageData.StartDate.Format("2006-01-02")
	stringMap["end_date"] = reservationPageData.EndDate.Format("2006-01-02")

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservationPageData"] = reservationPageData
	render.ParseTemplet(w, r, "reservation-summary.page.tmpl", &models.TempletData{
		Data:      data,
		StringMap: stringMap,
	})
}

// Renders the Choose Room page
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "roomid"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot find reservations data in the sessions"))
		return
	}

	res.RoomId = roomId
	room, err := m.DB.GetRoomById(roomId)
	if err != nil {
		helpers.ServerError(w, err)
	}
	res.Room = room
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)

}
