package main

import (
	"html/template"

	"golang.org/x/text/language"
)

const ServerListen = ":7000"

const EntriesDirectory = "./entries"
const TemplateFile = "./public/index.html"
const StaticDirectory = "./public/static"
var ContentLanguage = language.English

const MainTitle = "Encyclopedia"
var FooterContent = []template.HTML{
	template.HTML("powered by <a href='https://github.com/ChaoticByte/plaintext-encyclopedia' target='_blank' rel='noopener noreferrer'>plaintext-encyclopedia</a>"),
}
