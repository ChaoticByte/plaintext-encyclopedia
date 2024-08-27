package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
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
	entryName := strings.Trim(req.URL.Path, "/")
	if entryName != "" {
		if strings.Contains(entryName, "/") || strings.Contains(entryName, ".") {
			// path traversal
			logger.Println("Possible path traversal attempt from", req.RemoteAddr, "to", entryName)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// load entry
		entry = db.Entries[entryName]
	}
	err = appTemplate.ExecuteTemplate(
		w, "app",
		TemplateData{TOC: db.Keys, Entry: entry, Title: MainTitle, EntryTitle: entryName})
	if err != nil { logger.Println(err) }
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
	// start server
	logger.Println("Starting server on", ServerListen)
	err := http.ListenAndServe(ServerListen, nil)
	if err != nil { logger.Panicln(err) }
}
