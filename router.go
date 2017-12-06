package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/FidelityInternational/cf-loupe/applist"
	"github.com/FidelityInternational/cf-loupe/cf"
	"github.com/julienschmidt/httprouter"
)

// caches response App Data
type crAppData struct {
	marshalledAppData *[]byte
	lastFetched       time.Time
	activelyScraping  *bool
}

// BuildRouter returns the main router
func BuildRouter(cfClients map[string]cf.IClient, timeNow func() time.Time) *httprouter.Router {
	crAppData := crAppData{
		activelyScraping: new(bool),
	}

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		templ, err := template.ParseFiles("templates/index.html")
		if err != nil {
			renderInternalServerError(w, err)
			return
		}

		if err = templ.Execute(w, nil); err != nil {
			log.Println(err.Error())
			return
		}
	})

	router.GET("/listapps", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jAppData, err := crAppData.scrape(cfClients, timeNow)
		if err != nil {
			renderInternalServerError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jAppData)
	})

	return router
}

func renderInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func (crAppData *crAppData) scrape(cfClients map[string]cf.IClient, timeNow func() time.Time) ([]byte, error) {
	now := timeNow()
	if crAppData.lastFetched.Before(now.Add(-60*time.Second)) && !*crAppData.activelyScraping {
		crAppData.activelyScraping = setPointerBool(true)
		appData, err := applist.BuildAppData(cfClients, now)
		if err != nil {
			return nil, err
		}

		jAppData, err := json.Marshal(appData)
		if err != nil {
			return nil, err
		}

		crAppData.marshalledAppData = &jAppData
		crAppData.lastFetched = time.Now()

		crAppData.activelyScraping = setPointerBool(false)
	}

	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		if !*crAppData.activelyScraping {
			break
		}
	}

	return *crAppData.marshalledAppData, nil
}

func setPointerBool(b bool) *bool {
	return &b
}
