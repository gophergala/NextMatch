package main

import (
	"flag"
	"github.com/gophergala/NextMatch/updater/xmlstats"
	"log"
	"net/http"
	"strings"
	"time"
	"os"
)

var port = flag.String(`p`, `80`, `the port on wich we're serving`)

func init() {
	addTfunc(`parse`, time.Parse)
}

func handler(w http.ResponseWriter, r *http.Request) {

	e , err := xmlstats.BySport("nba", "20150123")
	if err != nil {
	    log.Printf("Didn't get data :( err = %v", err)
	} else {
	    log.Printf("Working with %#v\n\n", e)
	}

	renderArgs := args{
		`events`: e,
		`title`:  `Home`,
	}

	customT(w, `home`, renderArgs)
}

func reload(rw http.ResponseWriter, req *http.Request) {
	loadTmpl()
	http.Redirect(rw, req, `/`, 302)
}

func main() {

	xmlstats.Token = os.Getenv("XMLSTATS_TOKEN")
	if len(xmlstats.Token) == 0 {
		log.Fatal("Specify XMLSTATS_TOKEN environment variable")
	}

	flag.Parse()
	loadTmpl()
	http.HandleFunc(`/refresh`, reload)
	http.Handle(`/static/`, static(http.FileServer(http.Dir(`./web`))))
	http.HandleFunc(`/`, handler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request){})
	log.Fatal(http.ListenAndServe(`:`+*port, nil))
}

func static(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, `/static`, `/resources`, 1)
		h.ServeHTTP(w, r)
	}
}
