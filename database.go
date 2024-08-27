package main

import (
	"io/fs"
	"log"
	"os"
	"strings"
)

type Database struct {
	Keys []string
	Entries map[string]string
}

func BuildDB(directory string) Database {
	logger := log.Default()
	logger.Println("Building database")
	// keys, entries
	var keys []string
	entries := map[string]string{}
	// get files in directory and read them
	directory = strings.TrimRight(directory, "/") // we don't need that last / -> if '/' is used as directory, you are dumb.
	entriesDirFs := os.DirFS(directory)
	keys, err := fs.Glob(entriesDirFs, "*")
	if err != nil { logger.Panicln(err) }
	for _, k := range keys {
		contentB, err := os.ReadFile(directory + "/" + k)
		if err != nil { logger.Panicln(err) }
		entries[k] = string(contentB)
	}
	// create bleve index
	return Database{Keys: keys, Entries: entries}
}
