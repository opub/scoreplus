package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

func routeMembers(r *chi.Mux) {
	r.Route("/member", func(r chi.Router) {

		//profile
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			templateHandler("profile", providers, w, r)
		})

	})
}
