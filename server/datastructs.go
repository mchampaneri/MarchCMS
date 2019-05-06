package main

// SlingRoute is routing entity for SlingPages
type SlingRoute struct {
	PageURL    string `json:"PageURL"`
	PageNumber string `json:"PageNumber" storm:"id,increment"`
}

// SlingPage is root struct for SlignPages
type SlingPage struct {
	PageTitle  string           `json:"PageTitle"`
	PageNumber string           `json:"PageNumber" storm:"id,increment"`
	Content    SlingPageContent `json:"PageContent"`
}

// SlingPageContent holds content of the page
type SlingPageContent struct {
	Keywords    string `json:"Keywords"`
	Descritpion string `json:"Desc"`
	HTML        string `json:"HTML"`
}
