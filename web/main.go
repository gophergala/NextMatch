package main

import (
	"flag"
	"log"
	"net/http"
)

var port = flag.String(`p`, `80`, `the port on wich we're serving`)

func handler(w http.ResponseWriter, r *http.Request) {
	customT(w, `home`, struct{ Content string }{`Hello from go`})
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(`:`+*port, nil))
}
