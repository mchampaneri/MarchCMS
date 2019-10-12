package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"
)

// Config holds global configurations of cms
type Config struct {
	ID         string `json:"id"`
	Address    string `json:"Address"`
	Name       string `json:"Name"`
	Database   string `json:"Database"`
	Theme      string `json:"Theme"`
	Status     string `json:"Status"`
	Live       string `json:"Live"`
	FaviconURL string `json:"FaviconURL"`
	LogoURL    string `json:"LogoURL"`
}

// ThemeConfig holds configrations for particular theme [ 1 copy per theme]
type ThemeConfig struct {
	ID    string       `json:"ID" storm:"id,increment"`
	Theme string       `json:"Theme"`
	Menus []*ThemeMenu `json:"Menus"`
}

// ThemeMenu holds MarchCMS <-> ThemeMenu maps
type ThemeMenu struct {
	ID    int    `json:"ID" storm:"id,increment"`
	Place string `json:"Place"`
	Menu  string `json:"Menu"`
}

// MarchPage is root struct for SlignPages

type MarchPage struct {
	PageNumber   string           `json:"PageNumber" storm:"id"`
	PageTemplate string           `json:"PageTemplate"`
	PageURL      string           `json:"PageURL" storm:"unique"`
	PageTitle    string           `json:"PageTitle"`
	Content      MarchPageContent `json:"PageContent"`
	Co           time.Time        `json:"Co" storm:"index"`
	Uo           time.Time        `json:"Uo" storm:"index"`
	Do           time.Time        `json:"Do" storm:"index"`
	MarchUserID  int              `json:"UserID" storm:"index"`
	UpdaterID    int              `json:"UpdaterUserID" storm:"index"`
	MarchUserObj MarchUser        `json:"Author"`
}

// MarchPost is root struct for SlignPages
type MarchPost struct {
	PageNumber   string           `json:"PageNumber" storm:"id"`
	PageTemplate string           `json:"PageTemplate"`
	PageURL      string           `json:"PageURL" storm:"unique"`
	PageTitle    string           `json:"PageTitle"`
	PageThumb    string           `json:"PageThumb"`
	Content      MarchPageContent `json:"PageContent"`
	Co           time.Time        `json:"Co" storm:"index"`
	Uo           time.Time        `json:"Uo" storm:"index"`
	Do           time.Time        `json:"Do" storm:"index"`
	Tag1         string           `json:"Tag1" storm:"index"`
	Tag2         string           `json:"Tag2" storm:"index"`
	Tag3         string           `json:"Tag3" storm:"index"`
	Type         string           `json:"Type" storm:"index"`
	MarchUserID  int              `json:"UserID" storm:"index"`
	UpdaterID    int              `json:"UpdaterUserID" storm:"index"`
	MarchUserObj MarchUser        `json:"Author"`
}

// func (p *MarchPost) MarshalJSON() ([]byte, error) {
// 	type Alias MarchPost
// 	marchUser := MarchUser{}
// 	if err := db.One("ID", p.MarchUserID, &marchUser); err == nil {
// 		log.Println(marchUser)
// 		return json.Marshal(&struct {
// 			MarchUserObj MarchUser `json:"Author"`
// 			*Alias
// 		}{
// 			MarchUserObj: marchUser,
// 			Alias:        (*Alias)(p),
// 		})
// 	} else {
// 		log.Fatalln("colud not fetch user:", err.Error())
// 		return []byte("hello"), nil
// 	}

// }

// MarchPageContent holds content of the page
type MarchPageContent struct {
	Keywords string `json:"Keywords"`
	Desc     string `json:"Desc"`
	HTML     string `json:"HTML"`
}

// MarchMenu holds content for navigation menu
type MarchMenu struct {
	ID          int                   `json:"ID" storm:"id,increment"`
	Slug        string                `json:"Slug"`
	Name        string                `json:"Name" storm:"unique"`
	Items       []*MarchMenuItemIndex `json:"Items"`
	MarchUserID int                   `json:"UserID" storm:"index"`
	UpdaterID   int                   `json:"UpdaterUserID" storm:"index"`
	Co          time.Time             `json:"Co" storm:"index"`
	Uo          time.Time             `json:"Uo" storm:"index"`
	Do          time.Time             `json:"Do" storm:"index"`
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
	ID          int            `json:"ID" storm:"id,increment"` // primary key
	Name        string         `json:"Name"`
	Email       string         `json:"Email"`
	Picture     string         `json:"Picture"`
	SmallDesc   string         `json:"SmallDesc"`
	Password    string         `json:"Password"`
	Role        int            `json:"Role"`
	Status      int            `json:"Status"`
	Extra       MarchUserExtra `json:-`
	MarchUserID int            `json:"UserID" storm:"index"`
	UpdaterID   int            `json:"UpdaterUserID" storm:"index"`
}

// MarchUserExtra holds extra information about user
type MarchUserExtra struct {
	Website     string `json:"Website"`
	Email       string `json:"Email"`
	Achivements string `json:Achivements`
}

// RegisterUser registers incomming marchUser with CMS
func (u *MarchUser) RegisterUser() (*MarchUser, error) {
	u.Password = u.HashPassword(u.Password)
	err := db.Save(u)
	if err != nil {
		color.Red("error during savaing the model", err.Error())
		return nil, err
	}
	// SendMail("mails/regdone.html",
	// 	"info.safelms@puberstreet.com",
	// 	"Register with Safelms",
	// 	"You have succesfully registred with the safelms ",
	// 	u.Email,
	// 	u)
	return u, nil
}

func AllUsers() (userlist []MarchUser) {
	db.All(&userlist)
	return
}

func (*MarchUser) HashPassword(password string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashBytes)
}

// Core Functions //
func (u *MarchUser) LoginUser() (authPass bool, su MarchUser) {
	err := db.One("Email", u.Email, &su)
	if err != nil {
		color.Red("error at login %s", err.Error())
	}
	if err := bcrypt.CompareHashAndPassword([]byte(su.Password), []byte(u.Password)); err == nil {
		authPass = true
	} else {
		authPass = false
		log.Println("failed to match pass", err.Error(), su.Password, u.Password)
	}
	return
}

func issueSession(w http.ResponseWriter, r *http.Request, su MarchUser) bool {
	usession, err := UserSession.Get(r, "mvc-user-session")
	if err != nil {
		log.Println("failed to issue session :", err.Error())
	}
	usession.Values["id"] = su.ID
	usession.Values["name"] = su.Name
	usession.Values["email"] = su.Email
	usession.Values["auth"] = true
	usession.Values["role"] = su.Role
	usession.Values["picture"] = su.Picture
	if su.Status == activeAccount {
		usession.Values["active"] = true
		usession.Values["message"] = "Welcome"
	} else {
		usession.Values["active"] = false
		usession.Values["message"] = "Please Active Your Account by Verifying Your Email Address"
	}
	log.Println("Session Issued")
	usession.Save(r, w)
	log.Println("session issued for :", usession.Values["name"])
	renderJSON(w, map[string]string{"success": "authentication done"})
	return true
}
