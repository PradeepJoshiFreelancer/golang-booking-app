package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/pradeepj4u/bookings/cmd/models"
	"github.com/pradeepj4u/bookings/internal/config"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func SetDefaultData(td *models.TempletData, r *http.Request) *models.TempletData {
	td.InfoEdit = app.Session.PopString(r.Context(), "InfoEdit")
	td.WarningEdit = app.Session.PopString(r.Context(), "WarningEdit")
	td.CriticalEdit = app.Session.PopString(r.Context(), "CriticalEdit")

	td.CSRFToken = nosurf.Token(r)
	return td
}

var pathToTemplet = "./templet"

// Renders the read file to browser
func ParseTemplet(w http.ResponseWriter, r *http.Request, t string, templetData *models.TempletData) {
	var templetCache map[string]*template.Template
	if app.UseCache {
		//create templet cache
		templetCache = app.TempletCache
	} else {
		templetCache, _ = CreateChacheMap()
	}

	//get requested templet from cache
	temp, ok := templetCache[t]
	if !ok {
		log.Fatal(ok)
	}

	buf := new(bytes.Buffer)
	templetData = SetDefaultData(templetData, r)
	err := temp.Execute(buf, templetData)

	if err != nil {
		log.Println(err)
	}
	//render templet
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateChacheMap() (map[string]*template.Template, error) {
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
