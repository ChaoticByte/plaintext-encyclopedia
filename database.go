package main

import (
	"io/fs"
	"log"
	"os"
	"strings"
)

type Entry struct {
	Name string		// from filename
	Content string	// from content
}

type Database struct {
	Entries []Entry
}

func (db *Database) TOC() []string {
	var names []string
	for _, e := range db.Entries {
		names = append(names, e.Name)
	}
	return names
}

func (db *Database) Entry(name string) string {
	// returns empty string if not found
	result := ""
	for _, e := range db.Entries {
		if e.Name == name {
			result = e.Content
		}
	}
	return result
}

func BuildDB(directory string) Database {
	logger := log.Default()
	logger.Println("Building database")
	var entries []Entry
	// get files in directory and read them
	directory = strings.TrimRight(directory, "/") // we don't need that last / -> if '/' is used as directory, you are dumb.
	entriesDirFs := os.DirFS(directory)
	files, err := fs.Glob(entriesDirFs, "*")
	if err != nil { logger.Panicln(err) }
	for _, f := range files {
		contentB, err := os.ReadFile(directory + "/" + f)
		if err != nil { logger.Panicln(err) }
		entries = append(entries, Entry{Name: f, Content: string(contentB)})
	}
	return Database{Entries: entries}
}
