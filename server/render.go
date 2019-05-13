// generatign responses in various foramts
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/CloudyKit/jet"
	blackfriday "gopkg.in/russross/blackfriday.v2"
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
		var AllPages []SlingPage
		db.All(&AllPages)
		return reflect.ValueOf(AllPages)
	})

	////////////// Admin //////////////////
	adminInstance.AddGlobalFunc("SiteTitle", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Name)
	})

	adminInstance.AddGlobalFunc("PageTemplates", func(a jet.Arguments) reflect.Value {
		// Read Pages templates
		templates := make([]string, 0, 10)
		if fileInfo, err := ioutil.ReadDir(filepath.Join(root, "themes", "shortshot", "pages")); err == nil {
			for _, file := range fileInfo {
				templates = append(templates, file.Name())
			}
			return reflect.ValueOf(templates)
		}
		return reflect.ValueOf(nil)
	})

}

func renderPage(w io.Writer, page SlingPage) {
	var pageTemplate = "index.html"
	if page.PageTemplate != "" && page.PageTemplate != "-" {
		pageTemplate = filepath.Join(".", "pages", page.PageTemplate)
	}
	// log.Fatalln(filepath.Join(root, "themes", "shortshot", "pages", pageTemplate))
	t, err := frontInstance.GetTemplate(pageTemplate)
	dataMap := map[string]interface{}{
		"Page": page,
		// "Route": r,
	}
	log.Println(page)

	output := blackfriday.Run([]byte(page.Content.HTML))
	dataMap["output"] = output
	if err = t.Execute(w, nil, dataMap); err != nil {
		log.Fatalf(" - respnose-generator.go  View  : %s", err.Error())
	}
}

func renderPagev1(w io.Writer, r SlingRoute) {
	t, err := frontInstance.GetTemplate("index.html")
	var page SlingPage
	err = db.One("PageNumber", r.PageNumber, &page)
	dataMap := map[string]interface{}{
		"Page":  page,
		"Route": r,
	}
	log.Println(page)

	output := blackfriday.Run([]byte(page.Content.HTML))
	dataMap["output"] = output
	if err = t.Execute(w, nil, dataMap); err != nil {
		log.Fatalf(" - respnose-generator.go  View  : %s", err.Error())
	}
}

func renderAdmin(w io.Writer, page string, data map[string]interface{}) {
	log.Println("Render admin is executing")
	t, err := adminInstance.GetTemplate(page)
	// dataMap := map[string]interface{}{}
	if err = t.Execute(w, nil, data); err != nil {
		log.Fatalf(" - respnose-generator.go  View  : %s", err.Error())
	}
}

// JSON Returns the data in form of "JSON" for the incoming
// request
func renderJSON(w http.ResponseWriter, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to generate json ")
	}
	fmt.Fprint(w, string(response))
}
