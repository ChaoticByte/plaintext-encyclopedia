package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var logger *log.Logger
var appTemplate *template.Template = template.New("app")
var db Database

type TemplateData struct {
	TOC []string
	Entry string
	Title string
	EntryTitle string
}

func loadTemplate() {
	data, err := os.ReadFile(TemplateFile)
	if err != nil { logger.Panicln(err) }
	appTemplate.Parse(string(data))
}

func handleApplication(w http.ResponseWriter, req *http.Request) {
	var entry string
	var err error
	entryName := path.Base(req.URL.Path)
	if entryName != "/" {
		// load entry
		entry = db.Entries[entryName]
		if entry == "" { // redirect if entry doesn't exist (or is empty)
			http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		}
	}
	err = appTemplate.ExecuteTemplate(
		w, "app",
		TemplateData{TOC: db.Keys, Entry: entry, Title: MainTitle, EntryTitle: entryName})
	if err != nil { logger.Println(err) }
}

func handleSearchAPI(w http.ResponseWriter, req *http.Request) {
	searchQuery := path.Base(req.URL.Path)
	if searchQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results := db.search(searchQuery)
	w.Write([]byte(strings.Join(results, "\n")))
}

func main() {
	// get logger
	logger = log.Default()
	// load template
	loadTemplate()
	// build db
	db = BuildDB(EntriesDirectory)
	// handle static files
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(StaticDirectory)))
	http.Handle("/static/", staticHandler)
	// handle application
	http.HandleFunc("/", handleApplication)
	http.HandleFunc("/search/", handleSearchAPI)
	// start server
	logger.Println("Starting server on", ServerListen)
	err := http.ListenAndServe(ServerListen, nil)
	if err != nil { logger.Panicln(err) }
}
