package main

import (
	"log"
	"net/http"
)

/* HTTP route middlewares
| Works on request level
*/

/* Authenticates the user against the credentials stored
| in the user table in the database and sets the session
| as well as a session cookie of the gorilla session for the
| further operations that requires the  authenticated
| User only.
*/

func check(req *http.Request) bool {
	usession, err := UserSession.Get(req, "mvc-user-session")
	if (err != nil) || (usession.Values["auth"] != true) {
		if err != nil {
			log.Println("error in auth:", err.Error())
		}
		return false
	}
	return true
}

func auth(pass http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if check(r) == true {
			pass(w, r)

		} else {
			// redirect to login
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}
	}
}

func author(pass http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if check(r) == true {
			session, _ := UserSession.Get(r, "mvc-user-session")
			if (session.Values["role"].(int) == writerUser) ||
				(session.Values["role"].(int) == editorUser) ||
				(session.Values["role"].(int) == adminUser) {
				pass(w, r)

			}
		}
		http.Redirect(w, r, r.Referer(), http.StatusUnauthorized)
		// redirect to login
	}
}

/* Operational Middlewares
| middelware that operates on particular opertaions
*/

/* Same user middleware to verify that operation is being done
| by the same user has created the resource
*/
func originalWriter(originalWritersID int, r *http.Request) bool {
	if session, err := UserSession.Get(r, "mvc-user-session"); err == nil {
		if originalWritersID == session.Values["id"].(int) {
			return true
		}
	}
	return false
}
