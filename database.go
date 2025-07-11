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
	Titles map[string]string
	Entries map[string]string
	matcher *search.Matcher
}

func (db *Database) searchForIds(query string) []string { // returns keys (entry names)
	results := []string{}
	// compile patterns
	queryPatterns := []*search.Pattern{}
	for _, q := range strings.Split(query, " ") { // per word
		queryPatterns = append(queryPatterns, db.matcher.CompileString(q))
	}
	// search
	for k, v := range db.Entries {
		// title (k)
		titleLower := strings.ToLower(db.Titles[k])
		queryLower := strings.ToLower(query)
		if strings.Contains(titleLower, queryLower) {
			results = append(results, k)
			continue
		}
		// content body
		patternsFound := 0
		for _, p := range queryPatterns {
			if s, _ := p.IndexString(v); s != -1 {
				patternsFound++ // this pattern was found
			}
		}
		if patternsFound == len(queryPatterns) && !slices.Contains(results, k) {
			// if all patterns were found, add the key (entry name) to the list
			results = append(results, k)
		}
	}
	slices.Sort(results)
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
	titles := map[string]string{}
	for _, f := range files {
		k := f[:len(f)-4] // remove ".txt"
		k = strings.ReplaceAll(k, "|", "_") // we don't want | because it is used in the search protocol
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
			titles[k] = title
			entries[k] = body
		} else {
			titles[k] = k
			entries[k] = content
		}
	}
	matcher := search.New(ContentLanguage, search.IgnoreCase, search.IgnoreDiacritics)
	return Database{Titles: titles, Entries: entries, matcher: matcher}
}
