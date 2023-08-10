package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pradeepj4u/bookings/cmd/models"
	"github.com/pradeepj4u/bookings/internal/config"
	"github.com/pradeepj4u/bookings/internal/forms"
	"github.com/pradeepj4u/bookings/internal/helpers"
	"github.com/pradeepj4u/bookings/internal/render"
)

// Stores the Class address
var Repo *Repository

// New Class Repository
type Repository struct {
	App *config.AppConfig
}

// Creates a new instance of class Repository
func NewRpository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
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
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// Returns a JSON response
func (m *Repository) CheckAvailabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "This is the response Message",
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

// Check availability post
func (m *Repository) PostCheckAvailability(w http.ResponseWriter, r *http.Request) {
	startDate := r.Form.Get("start")
	endDate := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Posted on Check availability page start Date %s and endDate %s", startDate, endDate)))
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
	var emptyReservation models.FormsData
	data := make(map[string]interface{})

	data["makeReservationPageData"] = emptyReservation

	render.ParseTemplet(w, r, "make-reservations.page.tmpl", &models.TempletData{
		Form: forms.NewForm(nil),
		Data: data,
	})
}
func (m *Repository) PostMakeReservations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
	}
	makeReservationPageData := models.FormsData{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	newForm := forms.NewForm(r.PostForm)

	newForm.Required("first_name", "last_name", "email", "phone")
	newForm.MinLength("first_name", 3, r)
	newForm.IsEmailValid("email")

	if !newForm.IsFormValid() {
		data := make(map[string]interface{})
		data["makeReservationPageData"] = makeReservationPageData
		render.ParseTemplet(w, r, "make-reservations.page.tmpl", &models.TempletData{
			Form: newForm,
			Data: data,
		})
		return
	}
	m.App.Session.Put(r.Context(), "makeReservationPageData", makeReservationPageData)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Renders the About page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservationPageData, ok := m.App.Session.Get(r.Context(), "makeReservationPageData").(models.FormsData)
	if !ok {

		m.App.Session.Put(r.Context(), "CriticalEdit", "Can't get makeReservationPageData in Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		m.App.ErrorLog.Println("Unable to get data from Session")
		return
	}
	m.App.Session.Remove(r.Context(), "makeReservationPageData")

	data := make(map[string]interface{})
	data["reservationPageData"] = reservationPageData
	render.ParseTemplet(w, r, "reservation-summary.page.tmpl", &models.TempletData{
		Data: data,
	})
}
