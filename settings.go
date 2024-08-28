package main

import "golang.org/x/text/language"

const ServerListen = ":7000"
const EntriesDirectory = "./entries"
const TemplateFile = "./public/index.html"
const StaticDirectory = "./public/static"
const MainTitle = "Encyclopedia"
var ContentLanguage = language.English
