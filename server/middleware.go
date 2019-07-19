package main

import (
	"net/http"
)

/* Creating the Session Singleton that going
|  to used by the other parts of the app
*/

/* Authenticates the user against the credentials stored
| in the user table in the database and sets the session
| as well as a session cookie of the gorilla session for the
| further operations that requires the  authenticated
| User only.
*/

func check(req *http.Request) bool {
	session, err := UserSession.Get(req, "mvc-user-session")
	if (err != nil) || session.IsNew || (session.Values["auth"] == false) {
		return false
	}
	return true
}

func auth(pass http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if check(r) == true {
			pass(w, r)
			return
		}
		// redirect to login
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func author(pass http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if check(r) == true {
			session, _ := UserSession.Get(r, "mvc-user-session")
			if (session.Values["role"] == writerUser) ||
				(session.Values["role"] == editorUser) {
				pass(w, r)
				return
			}
		}
		http.Redirect(w, r, r.Referer(), http.StatusUnauthorized)
		// redirect to login
	}
}
