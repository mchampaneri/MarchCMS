package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func adminRoutes(router *mux.Router) {

	// dashboard
	router.HandleFunc("/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/dashboard.html", map[string]interface{}{})
	})

	// content routes
	router.HandleFunc("/admin/page/create", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/page-create.html", map[string]interface{}{})
	})

	router.HandleFunc("/admin/page/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/page-edit.html", map[string]interface{}{})
	})

	router.HandleFunc("/admin/pages/list", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/pages.html", map[string]interface{}{})
	})

	// settigns routes

	// main route
	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/index.html", map[string]interface{}{})
	})

}
