// generatign responses in various foramts
package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/CloudyKit/jet"
)

// Root directory for view files
// [ where html templates are stored .. ]
var root, _ = os.Getwd()
var frontInstance = jet.NewHTMLSet(filepath.Join(root, "themes", "shortshot"))
var adminInstance = jet.NewHTMLSet(filepath.Join(root, "admin"))

func init() {

	frontInstance.SetDevelopmentMode(true)
	adminInstance.SetDevelopmentMode(true)

	/////////////// Front ///////////////
	frontInstance.AddGlobalFunc("SiteTitle", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Name)
	})

	frontInstance.AddGlobalFunc("PageList", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Name)
	})

	// frontInstance.AddGlobalFunc("PageContent", func(a jet.Arguments) reflect.Value {
	// 	return reflect.ValueOf(config.Name)
	// })

	frontInstance.AddGlobalFunc("PageList", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(sampleRoutes)
	})

	////////////// Admin //////////////////
	adminInstance.AddGlobalFunc("SiteTitle", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Name)
	})

}

func renderPage(w io.Writer, r SlingRoute) {
	t, err := frontInstance.GetTemplate("index.html")
	var page SlingPage
	err = db.One("PageNumber", r.PageNumber, &page)
	dataMap := map[string]interface{}{
		"Page":  page,
		"Route": r,
	}
	log.Println(page)
	if err = t.Execute(w, nil, dataMap); err != nil {
		log.Fatalf(" - respnose-generator.go  View  : %s", err.Error())
	}
}

func renderAdmin(w io.Writer, page string, data map[string]interface{}) {
	t, err := adminInstance.GetTemplate(page)
	// dataMap := map[string]interface{}{}
	if err = t.Execute(w, nil, data); err != nil {
		log.Fatalf(" - respnose-generator.go  View  : %s", err.Error())
	}
}
