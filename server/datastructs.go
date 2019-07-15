package main

import "time"

// Config holds global configurations of cms
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
	ID    int                   `json:"ID" storm:"id,increment"` // primary key
	Slug  string                `json:"Slug"`
	Name  string                `json:"Name" storm:"unique"`
	Items []*MarchMenuItemIndex `json:"Items"`
}

// MarchMenuItemIndex holds content for navigation menu
type MarchMenuItemIndex struct {
	ID   int            `json:"ID" storm:"id,increment"` // primary key
	Item *MarchMenuItem `json:"Item"`
}

// MarchMenuItem holds individual menu item for menu
type MarchMenuItem struct {
	ID    int    `json:"ID" storm:"id,increment"` // primary key
	Slug  string `json:"-"`
	Title string `json:"Name"`
	URL   string `json:"URL"`
	// CSSClass  string `json:"CSSClass"`
	// ElementID string `json:"ElementID"`
}

// MarchUser holds user information and session details
type MarchUser struct {
	ID        int            `json:"ID" storm:"id,increment"` // primary key
	Name      string         `json:"Name"`
	SmallDesc string         `json:"SmallDesc"`
	Extra     MarchUserExtra `json:-`
}

// MarchUser extra holds extra information about user
type MarchUserExtra struct {
	Website     string `json:"Website"`
	Email       string `json:"Email"`
	Achivements string `json:Achivements`
}
