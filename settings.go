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
var SiteDescription = `A very simple and clean reader/browser?
for plaintext... documentation? wiki? chapters of a book? whatever.`
var FooterContent = []template.HTML{
	template.HTML("powered by <a href='https://github.com/ChaoticByte/plaintext-encyclopedia' target='_blank' rel='noopener noreferrer'>plaintext-encyclopedia</a>"),
}
