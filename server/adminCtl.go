package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func adminRoutes(router *mux.Router) {

	router.HandleFunc("/admin/pages/list", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/pages.html", map[string]interface{}{})
	})
	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/index.html", map[string]interface{}{})
	})

}
