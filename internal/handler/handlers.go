package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		return
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
		return
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
		return
	}
	res.Room = room
	res.RoomId = roomId

	layout := "2006-01-02"
	sd, err := time.Parse(layout, startDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
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
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailibilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
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

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "CriticalEdit", "Can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "CriticalEdit", "Cannot parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "CriticalEdit", "Cannot parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "CriticalEdit", "Invalid data")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot find reservations data in the sessions"))
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")
	reservation.StartDate = startDate
	reservation.EndDate = endDate
	reservation.Room.ID = roomId
	reservation.RoomId = roomId

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	newForm := forms.NewForm(r.PostForm)

	newForm.Required("first_name", "last_name", "email", "phone")
	newForm.MinLength("first_name", 3, r)
	newForm.IsEmailValid("email")

	if !newForm.IsFormValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.ParseTemplet(w, r, "make-reservations.page.tmpl", &models.TempletData{
			Form:      newForm,
			Data:      data,
			StringMap: stringMap,
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
		return
	}
	res.Room = room
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.ParseTemplet(w, r, "login.page.tmpl", &models.TempletData{
		Form: forms.NewForm(nil),
	})
}

// handels login request and authentcates user
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	form := forms.NewForm(r.PostForm)
	form.Required("email", "password")
	form.IsEmailValid("email")

	if !form.IsFormValid() {
		// Route to the login pages
		render.ParseTemplet(w, r, "login.page.tmpl", &models.TempletData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {

		log.Println(err)

		m.App.Session.Put(r.Context(), "CriticalEdit", "Invalid Password")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "InfoEdit", "Login Successfull!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout the PPT out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Handels request form Authenticated users for Admin dashboard
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.ParseTemplet(w, r, "admin-dashboard.page.tmpl", &models.TempletData{})
}

// Handels request for new Reservations
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		log.Println(err)
	}
	data := make(map[string]interface{})

	data["reservations"] = reservations
	render.ParseTemplet(w, r, "admin-new-reservations.page.tmpl", &models.TempletData{
		Data: data,
	})
}

// Handels request for all Reservations
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		log.Println(err)
	}
	data := make(map[string]interface{})

	data["reservations"] = reservations

	render.ParseTemplet(w, r, "admin-all-reservations.page.tmpl", &models.TempletData{
		Data: data,
	})
}

// Handels request form AdminReservationsCalendar
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}
	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	data := make(map[string]interface{})
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear

	data["now"] = now

	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	//get first and last date of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfCurrentMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfCurrentMonth := firstOfCurrentMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)

	intMap["days_in_month"] = lastOfCurrentMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["rooms"] = rooms

	for _, room := range rooms {

		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfCurrentMonth; !d.After(lastOfCurrentMonth); d = d.AddDate(0, 0, 1) {

			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		//get all the restrictions for a given month

		restrictions, err := m.DB.GetRestrictionsByRoomForDate(room.ID, firstOfCurrentMonth, lastOfCurrentMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		for _, r := range restrictions {

			if r.ReservationId > 0 {
				//its a reservation
				for d := r.StartDate; !d.After(r.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = r.ReservationId
				}
			} else {
				//its a restriction
				blockMap[r.StartDate.Format("2006-01-2")] = r.ID
			}
		}
		data[fmt.Sprintf("reservation_map_%d", room.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", room.ID)] = blockMap
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", room.ID), blockMap)

	}

	render.ParseTemplet(w, r, "admin-reservations-calendar.page.tmpl", &models.TempletData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

// Handels request form Show Single reservation on the Pages
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	parsedURL := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(parsedURL[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	source := parsedURL[3]
	stringMap := make(map[string]string)
	stringMap["source"] = source

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["month"] = month
	stringMap["year"] = year

	//Get data for reservation id form database
	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	dataMap := make(map[string]interface{})
	dataMap["reservation"] = res

	render.ParseTemplet(w, r, "admin-show-reservation.page.tmpl", &models.TempletData{
		StringMap: stringMap,
		Data:      dataMap,
		Form:      forms.NewForm(nil),
	})
}

// Handels Post request to update the Reservation Form
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	parsedURL := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(parsedURL[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	source := parsedURL[3]
	stringMap := make(map[string]string)
	stringMap["source"] = source

	//Get data for reservation id form database
	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "InfoEdit", "Changes Saved!")
	//Redirect to all reservation page
	if year == "" {

		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", source), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

// Marks a Reservation processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = m.DB.UpdateProcessed(id, 1)

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "InfoEdit", "Reservation Processed!")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

// Marks a Reservation processed
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = m.DB.DeleteReservation(id)

	m.App.Session.Put(r.Context(), "InfoEdit", "Reservation Deleted!")
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

// Handels request form PostAdminReservationsCalendar
func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, err := strconv.Atoi(r.Form.Get("y"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	month, err := strconv.Atoi(r.Form.Get("m"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.NewForm(r.Form)
	for _, room := range rooms {
		// Get the blockMap form sessions(data before changes), if there is an entry in the blockMap that
		// does not exist in the posted data and restriction id > 0. We need to remove that restrictions.
		currBlockMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.ID)).(map[string]int)
		for name, value := range currBlockMap {
			if val, ok := currBlockMap[name]; ok {
				//only check for blocks that have valuws > 0 and that are not in form post
				// others are just days without blocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", room.ID, name)) {
						err := m.DB.DeleteRoomRestrictionById(value)
						if err != nil {
							helpers.ServerError(w, err)
							return
						}
					}
				}
			}
		}
	}

	//Logic to add a new Block
	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			parsedName := strings.Split(name, "_")
			roomId, _ := strconv.Atoi(parsedName[2])
			t, _ := time.Parse("2006-01-2", parsedName[3])

			//insert the new block
			err := m.DB.InsertRoomRestrictionForRoom(roomId, t)
			if err != nil {
				helpers.ServerError(w, err)
				return
			}

		}

	}
	m.App.Session.Put(r.Context(), "InfoEdit", "Changes saved!")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}
