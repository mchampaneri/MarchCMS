// generatign responses in various foramts
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/CloudyKit/jet"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var frontInstance = jet.NewHTMLSet(themesFolder)
var adminInstance = jet.NewHTMLSet(adminFolder)

func init() {

	frontInstance.SetDevelopmentMode(true)
	adminInstance.SetDevelopmentMode(true)

	/////////////// Front ///////////////
	frontInstance.AddGlobalFunc("SiteTitle", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Name)
	})

	frontInstance.AddGlobalFunc("PageList", func(a jet.Arguments) reflect.Value {
		var AllPages []MarchPage
		db.All(&AllPages)
		return reflect.ValueOf(AllPages)
	})

	////////////// Admin //////////////////
	adminInstance.AddGlobalFunc("SiteTitle", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Name)
	})

	adminInstance.AddGlobalFunc("ActiveTheme", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Theme)
	})

	// adminInstance.AddGlobalFunc("PluginMenu", func(a jet.Arguments) reflect.Value {

	// 	menus := make([]string, 1, 10)
	// 	for _, extension := range extensions {
	// 		request := Request{}
	// 		response := new(Response)
	// 		if err := extension.Call("Admin.HookInMenu", request, response); err == nil {
	// 			menus = append(menus, response.Output)
	// 		}
	// 	}
	// 	return reflect.ValueOf(menus)
	// })

	adminInstance.AddGlobalFunc("PageTemplates", func(a jet.Arguments) reflect.Value {
		// Read Pages templates
		templates := make([]string, 0, 10)
		if fileInfo, err := ioutil.ReadDir(filepath.Join(root, "themes", config.Theme, "pages")); err == nil {
			for _, file := range fileInfo {
				templates = append(templates, file.Name())
			}
			return reflect.ValueOf(templates)
		}
		return reflect.ValueOf(nil)
	})

	adminInstance.AddGlobalFunc("PostTemplates", func(a jet.Arguments) reflect.Value {
		// Read Pages templates
		templates := make([]string, 0, 10)
		if fileInfo, err := ioutil.ReadDir(filepath.Join(root, "themes", config.Theme, "posts")); err == nil {
			for _, file := range fileInfo {
				templates = append(templates, file.Name())
			}
			return reflect.ValueOf(templates)
		}
		return reflect.ValueOf(nil)
	})

	adminInstance.AddGlobalFunc("InstalledThemes", func(a jet.Arguments) reflect.Value {
		// Read Pages templates
		templates := make([]struct{ Name, Thumb string }, 0, 10)
		if fileInfo, err := ioutil.ReadDir(filepath.Join(root, "themes")); err == nil {
			for _, file := range fileInfo {
				templates = append(templates, struct{ Name, Thumb string }{file.Name(),
					fmt.Sprint("/admin/themes-thumb/", file.Name(), "/thumb.png")})
			}
			return reflect.ValueOf(templates)
		}
		return reflect.ValueOf(nil)
	})

	// adminInstance.AddGlobalFunc("AvailableExtensions", func(a jet.Arguments) reflect.Value {
	// 	// Read Pages templates
	// 	templates := make([]RpcExtension, 0, 10)
	// 	if fileInfo, err := ioutil.ReadDir(filepath.Join(root, "extensions")); err == nil {
	// 		for _, file := range fileInfo {
	// 			configFile := filepath.Join(extensionFolder, file.Name(), "config.json")
	// 			readFile, err := os.Open(configFile)
	// 			if err != nil {
	// 				log.Fatalln(err.Error())
	// 			}
	// 			var extensionConfig RpcExtension
	// 			configDecoder := json.NewDecoder(readFile)
	// 			configDecoder.Decode(&extensionConfig)
	// 			templates = append(templates, extensionConfig)
	// 		}
	// 		return reflect.ValueOf(templates)
	// 	}
	// 	return reflect.ValueOf(nil)
	// })

}

func renderPage(w io.Writer, page MarchPage) {
	var pageTemplate = "index.html"
	if page.PageTemplate != "" && page.PageTemplate != "-" {
		pageTemplate = filepath.Join(".", "pages", page.PageTemplate)
	}
	t, err := frontInstance.GetTemplate(filepath.Join(config.Theme, pageTemplate))
	if err != nil {
		log.Fatalln(pageTemplate, " - ", config.Theme, " - ", err.Error())
	}

	dataMap := map[string]interface{}{
		"Page": page,
	}
	log.Println(page)

	output := blackfriday.Run([]byte(page.Content.HTML))
	dataMap["output"] = output
	if err = t.Execute(w, nil, dataMap); err != nil {
		log.Fatalf(" - respnose-generator.go  View  : %s", err.Error())
	}
}

func renderPost(w io.Writer, post MarchPost) {
	var pageTemplate = "index.html"
	if post.PageTemplate != "" && post.PageTemplate != "-" {
		pageTemplate = filepath.Join(".", "posts", post.PageTemplate)
	}
	t, err := frontInstance.GetTemplate(filepath.Join(config.Theme, pageTemplate))
	if err != nil {
		log.Fatalln(pageTemplate, " - ", config.Theme, " - ", err.Error())
	}
	dataMap := map[string]interface{}{
		"Page": post,
	}
	output := blackfriday.Run([]byte(post.Content.HTML))
	dataMap["output"] = output
	if err = t.Execute(w, nil, dataMap); err != nil {
		log.Fatalf(" - respnose-generator.go  View  : %s", err.Error())
	}
}

func renderAdmin(w io.Writer, page string, data map[string]interface{}) {
	log.Println("Render admin is executing")
	// log.Fatalln(page)
	if t, err := adminInstance.GetTemplate(page); err == nil {
		// dataMap := map[string]interface{}{}
		if err := t.Execute(w, nil, data); err != nil {
			log.Fatalf(" - respnose-generator.go  View  : %s", err.Error())
		}
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
