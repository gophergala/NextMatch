package main

import (
	"flag"
	"github.com/gophergala/NextMatch/updater"
	"log"
	"net/http"
	"strings"
	"time"
)

var port = flag.String(`p`, `80`, `the port on wich we're serving`)

func init() {
	addTfunc(`parse`, time.Parse)
}

func handler(w http.ResponseWriter, r *http.Request) {
	e := updater.Events{}
	updater.Unmarshal(sampleData, &e)

	renderArgs := args{
		`events`: e,
		`title`:  `Home`,
	}

	execT(w, `home`, renderArgs)
}

func reload(rw http.ResponseWriter, req *http.Request) {
	loadTmpl()
	http.Redirect(rw, req, `/`, 302)
}

func main() {
	flag.Parse()
	loadTmpl()
	http.HandleFunc(`/refresh`, reload)
	http.Handle(`/static/`, static(http.FileServer(http.Dir(`./web`))))
	http.HandleFunc(`/`, handler)
	log.Fatal(http.ListenAndServe(`:`+*port, nil))
}

func static(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, `/static`, `/resources`, 1)
		h.ServeHTTP(w, r)
	}
}
