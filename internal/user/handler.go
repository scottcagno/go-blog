package user

import (
	"net/http"
)

func (u *UserService) LoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		u.t.Render(w, r, "login.html", nil)
	})
}

func (u *UserService) LogoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		u.t.Render(w, r, "logout.html", nil)
	})
}

func (u *UserService) UserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		u.t.Render(w, r, "user.html", nil)
	})
}

func (u *UserService) HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		u.t.Render(w, r, "home.html", nil)
	})
}
