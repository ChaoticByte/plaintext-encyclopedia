package main

// Copyright (c) 2024 Julian Müller (ChaoticByte)

import (
	"io/fs"
	"log"
	"os"
	"slices"
	"strings"

	"golang.org/x/text/search"
)

type Database struct {
	Keys []string
	Entries map[string]string
	matcher *search.Matcher
}

func (db *Database) search(query string) []string { // returns keys (entry names)
	results := []string{}
	// compile patterns
	queryPatterns := []*search.Pattern{}
	for _, q := range strings.Split(query, " ") { // per word
		queryPatterns = append(queryPatterns, db.matcher.CompileString(q))
	}
	// search
	for _, k := range db.Keys {
		// title (k)
		if strings.Contains(k, query) || strings.Contains(query, k) {
			results = append(results, k)
			continue
		}
		// content body
		patternsFound := 0
		for _, p := range queryPatterns {
			if s, _ := p.IndexString(db.Entries[k]); s != -1 {
				patternsFound++ // this pattern was found
			}
		}
		if patternsFound == len(queryPatterns) && !slices.Contains(results, k) {
			// if all patterns were found, add the key (entry name) to the list
			results = append(results, k)
		}
	}
	return results
}

func BuildDB(directory string) Database {
	logger := log.Default()
	logger.Println("Building database")
	// files, entries
	var files []string
	entries := map[string]string{}
	// get files in directory and read them
	directory = strings.TrimRight(directory, "/") // we don't need that last /, don't use the root directory /
	entriesDirFs := os.DirFS(directory)
	files, err := fs.Glob(entriesDirFs, "*.txt")
	if err != nil { logger.Panicln(err) }
	titles := []string{}
	for _, f := range files {
		fileData, err := os.ReadFile(directory + "/" + f)
		if err != nil { logger.Panicln(err) }
		content := string(fileData)
		content = strings.Trim(content, "\n\r ")
		if strings.HasPrefix(content, "Title:") {
			len_content := len(content)
			title_start_idx := 6
			var title_stop_idx int
			if strings.Contains(content, "\n") {
				title_stop_idx = strings.Index(content, "\n")
			} else {
				title_stop_idx = len_content
			}
			title := content[title_start_idx:title_stop_idx]
			title = strings.Trim(title, " ")
			body := content[title_stop_idx:len_content]
			body = strings.TrimLeft(body, "\n\r")
			if len(body) < 1 {
				body = " "
			}
			titles = append(titles, title)
			entries[title] = body
		} else {
			title := f[:len(f)-4] // remove ".txt"
			titles = append(titles, title)
			entries[title] = content
		}
	}
	matcher := search.New(ContentLanguage, search.IgnoreCase, search.IgnoreDiacritics)
	return Database{Keys: titles, Entries: entries, matcher: matcher}
}
