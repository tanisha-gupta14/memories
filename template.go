package main

import (
	"html/template"
	"log"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func init() {
	if templates == nil {
		log.Fatal("Failed to parse templates")
	}
}
