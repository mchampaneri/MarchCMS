package main

// SlingRoute is routing entity for SlingPages
type SlingRoute struct {
	PageURL    string `json:"page-url"`
	PageNumber string `json:"page-number" storm:"id,increment"`
	PageTitle  string `json:"page-title"`
}

// SlingPage is root struct for SlignPages
type SlingPage struct {
	PageNumber string           `json:"page-number" storm:"id,increment"`
	Content    SlingPageContent `json:"page-content"`
}

// SlingPageContent holds content of the page
type SlingPageContent struct {
	Keywords    string `json:"keywords"`
	Descritpion string `json:"description"`
	HTML        string `json:"html"`
}
