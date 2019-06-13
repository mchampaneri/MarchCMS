package main

import "time"

type RpcExtension struct {
	Name    string `json:"Name"`
	Address string `json:"Address"`
	Status  string `json:"Status"`
	Author  string `json:"Author"`
	Desc    string `json:"Desc"`
	Website string `json:"Website"`
	Licence string `json:"Licence"`
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

// MarchPage is root struct for SlignPages
type MarchPage struct {
	PageTemplate string           `json:"PageTemplate"`
	PageURL      string           `json:"PageURL" storm:"unique"`
	PageTitle    string           `json:"PageTitle"`
	PageNumber   string           `json:"PageNumber" storm:"id"`
	Content      MarchPageContent `json:"PageContent"`
	Co           time.Time        `json:"Co" storm:"index"`
	Uo           time.Time        `json:"Uo" storm:"index"`
	Do           time.Time        `json:"Do" storm:"index"`
}

// MarchPost is root struct for SlignPages
type MarchPost struct {
	PageTemplate string           `json:"PageTemplate"`
	PageURL      string           `json:"PageURL" storm:"unique"`
	PageTitle    string           `json:"PageTitle"`
	PageNumber   string           `json:"PageNumber" storm:"id"`
	Content      MarchPageContent `json:"PageContent"`
	Co           time.Time        `json:"Co" storm:"index"`
	Uo           time.Time        `json:"Uo" storm:"index"`
	Do           time.Time        `json:"Do" storm:"index"`
}

// MarchPageContent holds content of the page
type MarchPageContent struct {
	Keywords string `json:"Keywords"`
	Desc     string `json:"Desc"`
	HTML     string `json:"HTML"`
}

// MarchMenu holds content for navigation menu
type MarchMenu struct {
	Index uint64          `json:"Index" storm:"id,increment"`
	Name  string          `json:"Name" storm:"unique"`
	Items []MarchMenuItem `json:"Items`
}

// MarchMenuItem holds individual menu item for menu
type MarchMenuItem struct {
	Title     string `json:"Title"`
	URL       string `json:"URL"`
	CSSClass  string `json:"CSSClass"`
	ElementID string `json:"ElementID"`
}
