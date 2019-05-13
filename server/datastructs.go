package main

import "time"

// SlingRoute is routing entity for SlingPages
type SlingRoute struct {
	PageURL    string `json:"PageURL"`
	PageNumber string `json:"PageNumber" storm:"id"`
}

// SlingPage is root struct for SlignPages
type SlingPage struct {
	PageTemplate string           `json:"PageTemplate`
	PageURL      string           `json:"PageURL" storm:"unique"`
	PageTitle    string           `json:"PageTitle"`
	PageNumber   string           `json:"PageNumber" storm:"id"`
	Content      SlingPageContent `json:"PageContent"`
	Co           time.Time        `json:"Co" storm:"index"`
	Uo           time.Time        `json:"Uo" storm:"index"`
	Do           time.Time        `json:"Do" storm:"index"`
}

// SlingPageContent holds content of the page
type SlingPageContent struct {
	Keywords string `json:"Keywords"`
	Desc     string `json:"Desc"`
	HTML     string `json:"HTML"`
}
