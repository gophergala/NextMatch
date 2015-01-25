package main

import (
	"flag"
	"github.com/gophergala/NextMatch/updater/xmlstats"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	shortf = "20060102"
	lonfg  = "2006-01-2T15:04:05-07:00"
)

var port = flag.String("p", "80", "the port on wich we're serving")

func init() {
	addTfunc("parse", time.Parse)
	addTfunc("now", time.Now)
}

func handler(w http.ResponseWriter, r *http.Request) {

	e, err := xmlstats.BySport("nba")
	if err != nil {
		log.Printf("Didn't get data :( err = %v", err)
	}

	renderArgs := args{
		"events": e,
		"title":  "Home",
	}

	execT(w, "home", renderArgs)
}

func reload(rw http.ResponseWriter, req *http.Request) {
	loadTmpl()
	http.Redirect(rw, req, "/", 302)
}

func main() {

	xmlstats.Token = os.Getenv("XMLSTATS_TOKEN")
	if len(xmlstats.Token) == 0 {
		log.Fatal("Specify XMLSTATS_TOKEN environment variable")
	}

	flag.Parse()
	loadTmpl()
	r := mux.NewRouter()
	r.HandleFunc(`/sport/{name}`, sportHandle)
	r.HandleFunc("/refresh", reload)
	r.HandleFunc("/", handler)
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.Handle("/static/", static(http.FileServer(http.Dir("."))))
	http.Handle(`/`, r)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func static(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, "/static", "/resources", 1)
		h.ServeHTTP(w, r)
	}
}

func sportHandle(w http.ResponseWriter, req *http.Request) {
	e, err := getSport(req)
	if err != nil {
		log.Printf("Didn't get data :( err = %v", err)
	}

	renderArgs := args{"events": e}
	log.Print(e.Event)

	execT(w, "events", renderArgs)
}

func getSport(req *http.Request) (xmlstats.Events, error) {
	vars := mux.Vars(req)

	name := vars[`name`]
	date, ok := vars[`date`]
	if !ok {
		date = time.Now().Format(shortf)
	}

	return xmlstats.BySport(name, date)
}
