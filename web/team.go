package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/opub/scoreplus/model"
	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

//TeamData has data required for template params
type TeamData struct {
	Sports  []model.Sport
	Results interface{}
}

func routeTeams(r *chi.Mux) {
	r.Route("/team", func(r chi.Router) {
		r.Use(SessionCtx)

		//list
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			templateHandler("team/list", "", false, TeamData{Sports: model.Sports}, w, r)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			message := ""
			search := r.Form.Get("search")
			sport := r.Form.Get("sport")
			results, err := model.SearchTeams(search, sport)
			if err != nil {
				log.Error().Err(err).Msg("team search failed")
				message = fmt.Sprintf("Search failed: %s", err.Error())
			}

			templateHandler("team/list", message, len(message) == 0, TeamData{Sports: model.Sports, Results: results}, w, r)
		})

		//new
		r.Get("/new", func(w http.ResponseWriter, r *http.Request) {
			templateHandler("team/new", "", true, TeamData{Sports: model.Sports}, w, r)
		})

		r.Post("/new", func(w http.ResponseWriter, r *http.Request) {
			m := r.Context().Value(sessionKey).(*model.Member)

			r.ParseForm()

			v := strings.TrimSpace(r.Form.Get("venue"))
			venue := model.Venue{Name: v}
			if len(v) > 0 {
				err := venue.Save()
				if err != nil {
					log.Error().Err(err).Msg("venue save failed")
				}
			}

			t := model.Team{}
			t.Sport = model.Sport(r.Form.Get("sport"))
			t.Name = strings.TrimSpace(r.Form.Get("teamname"))
			t.Mascot = strings.TrimSpace(r.Form.Get("mascot"))
			t.Venue = venue
			t.CreatedBy = m.ID

			err := t.Save()
			if err != nil {
				log.Error().Err(err).Msg("new team failed")
				message := fmt.Sprintf("Update failed: %s", err.Error())

				templateHandler("team/new", message, false, TeamData{Sports: model.Sports, Results: t}, w, r)
			} else {
				path := fmt.Sprintf("/team/%s", t.LinkID())
				http.Redirect(w, r, path, http.StatusSeeOther)
			}
		})

		//details
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			s := chi.URLParam(r, "id")
			id := util.DecodeLink(s)
			m, err := model.GetTeamFull(id)
			if err != nil {
				log.Warn().Str("id", s).Msg("team not found")
				render.Render(w, r, ErrBadRequest(err))
				return
			}
			if m.ID == 0 {
				render.Render(w, r, ErrNotFound)
				return
			}
			templateHandler("team/details", "", true, m, w, r)
		})
	})
}
