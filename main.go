package main

import (
	"html/template"
	"net/http"
	"os"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("<h1>Hello World!</h1>"))
	tpl.Execute(w, nil)
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
	}
	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
