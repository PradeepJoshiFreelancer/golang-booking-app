package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pradeepj4u/bookings/cmd/models"
	"github.com/pradeepj4u/bookings/internal/config"
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
	stringMap["test"] = "This is the first data passes"
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.ParseTemplet(w, r, "home.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}

// Renders the About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{}
	stringMap["remote_ip"] = remoteIp

	render.ParseTemplet(w, r, "about.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}

// Renders the Room page
func (m *Repository) StandardRoom(w http.ResponseWriter, r *http.Request) {

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{}
	stringMap["remote_ip"] = remoteIp

	render.ParseTemplet(w, r, "standard-room.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}

// Renders the Check availability post
func (m *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{}
	stringMap["remote_ip"] = remoteIp

	render.ParseTemplet(w, r, "check-availability.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}

type jsonResponse struct {
	OK      bool
	Message string
}

// Returns a JSON response
func (m *Repository) CheckAvailabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "This is the response Message",
	}
	out, err := json.MarshalIndent(resp, "", "      ")
	if err != nil {
		log.Println(err)
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

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{}
	stringMap["remote_ip"] = remoteIp

	render.ParseTemplet(w, r, "king-suit.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}

// Renders the Contact Us page
func (m *Repository) CantactUs(w http.ResponseWriter, r *http.Request) {

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{}
	stringMap["remote_ip"] = remoteIp

	render.ParseTemplet(w, r, "contact.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}

// Renders the Check availability post
func (m *Repository) MakeReservations(w http.ResponseWriter, r *http.Request) {

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{}
	stringMap["remote_ip"] = remoteIp

	render.ParseTemplet(w, r, "make-reservations.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}
