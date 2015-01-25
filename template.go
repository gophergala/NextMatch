package main

import (
	"html/template"
	"net/http"
)

var (
	tLoader *template.Template
	tp      = `web/resources/templates/`
)

func init() {
	// cache template parsing
	templates := tmplPath(`default.html`, `home.html`)
	tLoader = template.Must(template.ParseFiles(templates...))
}

func execT(w http.ResponseWriter, name string, data interface{}) {
	if err := tLoader.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func head(w http.ResponseWriter, data interface{}) {
	execT(w, `header`, data)
}

func footer(w http.ResponseWriter, data interface{}) {
	execT(w, `footer`, data)
}

func content(w http.ResponseWriter, data interface{}) {
	execT(w, `content`, data)
}

func defaultT(w http.ResponseWriter, data interface{}) {
	head(w, data)
	content(w, data)
	footer(w, data)
}

func customT(w http.ResponseWriter, name string, data interface{}) {
	head(w, data)
	execT(w, name, data)
	footer(w, data)
}

func tmplPath(n ...string) []string {
	for i, v := range n {
		n[i] = tp + v
	}

	return n
}
