package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/pradeepj4u/bookings/internal/config"
)

var app *config.AppConfig

func NewHelper(a *config.AppConfig) {
	app = a
}

func ClientErrors(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client Error", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

func IsAuntenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "user_id")
}
