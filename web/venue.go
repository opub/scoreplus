package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/opub/scoreplus/model"
	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
)

func routeVenues(r *chi.Mux) {
	r.Route("/venue", func(r chi.Router) {
		r.Use(SessionCtx)

		//list
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			templateHandler("venue/list", "", false, nil, w, r)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			message := ""
			search := r.Form.Get("search")
			results, err := model.SearchVenues(search)
			if err != nil {
				log.Error().Err(err).Msg("venue search failed")
				message = fmt.Sprintf("Search failed: %s", err.Error())
			}

			templateHandler("venue/list", message, len(message) == 0, results, w, r)
		})

		//new
		r.Get("/new", func(w http.ResponseWriter, r *http.Request) {
			templateHandler("venue/new", "", true, nil, w, r)
		})

		r.Post("/new", func(w http.ResponseWriter, r *http.Request) {
			m := r.Context().Value(sessionKey).(*model.Member)

			r.ParseForm()

			v := model.Venue{}
			v.Name = strings.TrimSpace(r.Form.Get("venuename"))
			v.Address = strings.TrimSpace(r.Form.Get("address"))
			v.Coordinates = strings.TrimSpace(r.Form.Get("coordinates"))
			v.CreatedBy = m.ID

			err := v.Save()
			if err != nil {
				log.Error().Err(err).Msg("new venue failed")
				message := fmt.Sprintf("Creation failed: %s", err.Error())

				templateHandler("venue/new", message, false, v, w, r)
			} else {
				path := fmt.Sprintf("/venue/%s", v.LinkID())
				http.Redirect(w, r, path, http.StatusSeeOther)
			}
		})

		//details
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			s := chi.URLParam(r, "id")
			id := util.DecodeLink(s)
			m, err := model.GetVenue(id)
			if err != nil {
				log.Warn().Str("id", s).Msg("venue not found")
				renderBadRequest(err, w, r)
				return
			}
			if m.ID == 0 {
				renderNotFound(w, r)
				return
			}
			templateHandler("venue/details", "", true, m, w, r)
		})
	})
}
