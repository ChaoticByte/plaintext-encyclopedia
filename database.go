package main

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
	// keys, entries
	var keys []string
	entries := map[string]string{}
	// get files in directory and read them
	directory = strings.TrimRight(directory, "/") // we don't need that last /, don't use the root directory /
	entriesDirFs := os.DirFS(directory)
	keys, err := fs.Glob(entriesDirFs, "*")
	if err != nil { logger.Panicln(err) }
	for _, k := range keys {
		contentB, err := os.ReadFile(directory + "/" + k)
		if err != nil { logger.Panicln(err) }
		content := string(contentB)
		entries[k] = content
	}
	matcher := search.New(ContentLanguage, search.IgnoreCase, search.IgnoreDiacritics)
	return Database{Keys: keys, Entries: entries, matcher: matcher}
}
