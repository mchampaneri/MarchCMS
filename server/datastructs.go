package main

import "time"

type RpcExtension struct {
	Name    string `json:"Name"`
	Address string `json:"Address"`
	Status  string `json:"Status"`
}

type Response struct {
	Output string
	Type   string
	Status string
}

type Request struct {
	Type  string
	Input map[string]interface{}
}

// Config holds global configurations
// of cms
type Config struct {
	ID       string `json:"id"`
	Address  string `json:"Address"`
	Name     string `json:"Name"`
	Database string `json:"Database"`
	Theme    string `json:"Theme"`
	Status   string `json:"Status"`
}

// SlingPage is root struct for SlignPages
type SlingPage struct {
	PageTemplate string           `json:"PageTemplate"`
	PageURL      string           `json:"PageURL" storm:"unique"`
	PageTitle    string           `json:"PageTitle"`
	PageNumber   string           `json:"PageNumber" storm:"id"`
	Content      SlingPageContent `json:"PageContent"`
	Co           time.Time        `json:"Co" storm:"index"`
	Uo           time.Time        `json:"Uo" storm:"index"`
	Do           time.Time        `json:"Do" storm:"index"`
}

// SlingPost is root struct for SlignPages
type SlingPost struct {
	PageTemplate string           `json:"PageTemplate"`
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
