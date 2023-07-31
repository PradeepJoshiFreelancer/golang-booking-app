package handler

import (
	"net/http"

	"github.com/pradeepj4u/bookings/cmd/models"
	"github.com/pradeepj4u/bookings/pkg/config"
	"github.com/pradeepj4u/bookings/pkg/render"
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

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	stringMap["test"] = "This is the first data passes"
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.ParseTemplet(w, "home.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{}
	stringMap["remote_ip"] = remoteIp

	render.ParseTemplet(w, "about.page.tmpl", &models.TempletData{
		StringMap: stringMap,
	})
}
