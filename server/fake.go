package main

import (
	"log"

	"github.com/satori/go.uuid"
)

var sampleRoutes = []SlingRoute{
	SlingRoute{
		PageURL: "/home",
		// PageNumber: "xyz-1",

	},
	SlingRoute{
		PageURL: "/about",
		// PageNumber: "xyz-2",

	},
}

var samplePages = []SlingPage{
	SlingPage{
		PageTitle: "Home",
		// PageNumber: "xyz-1",
		Content: SlingPageContent{
			Desc:     "Home page is alway home to website",
			Keywords: "SlingPages, High Performace CMS, Compiled to navtive ",
			HTML:     "<h1>Home Page</h1>",
		},
	},
	SlingPage{
		// PageNumber: "xyz-2",
		PageTitle: "About",
		Content: SlingPageContent{
			Desc:     "Home page is alway home to website",
			Keywords: "SlingPages, High Performace CMS, Compiled to navtive ",
			HTML:     "<h2> This is about page</h2>",
		},
	},
}

func feedFakeData() {
	for i, route := range sampleRoutes {
		if uuid, err := uuid.NewV1(); err != nil {
			log.Fatalln("Failed to gerate page id :", err.Error())
		} else {
			route.PageNumber = uuid.String()
			if err := db.Save(&route); err != nil {
				log.Fatalln("Failed to save route :", err.Error())
			} else {
				samplePages[i].PageNumber = route.PageNumber
				if err := db.Save(&samplePages[i]); err != nil {
					log.Fatalln("failed to save page : ", err.Error())
				} else {
					log.Println("Route & page savde with id ", route.PageNumber)
				}
			}
		}
	}
}
