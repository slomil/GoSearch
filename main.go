package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/freshman-tech/news-demo-starter-files/news"
)

var tpl = template.Must(template.ParseFiles("index.html"))
var newsapi *news.Client

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("<h1>Hello World!</h1>"))

	//tpl.Execute(w, nil)

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}
		fmt.Println("Search query is:", searchQuery)
		//	fmt.Println("Page is:", page)

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:    results,
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)

		fmt.Printf("%+v", results)
	}
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Port is empty")
		port = "3005"
	}

	/*apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("API key is empty")
	}*/

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, "6c718b169bf1469c93829e1f26662b26", 20)

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/search", searchHandler(newsapi))
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
