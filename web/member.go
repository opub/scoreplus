package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/opub/scoreplus/model"
	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
)

var sessionKey = &contextKey{"Session"}

//MemberData has data required for template params
type MemberData struct {
	Follows bool
	Results interface{}
}

func routeMembers(r *chi.Mux) {
	r.Route("/member", func(r chi.Router) {
		r.Use(SessionCtx)

		//list
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			templateHandler("member/list", "", false, nil, w, r)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			message := ""
			search := r.Form.Get("search")
			results, err := model.SearchMembers(search)
			if err != nil {
				log.Error().Err(err).Msg("member search failed")
				message = fmt.Sprintf("Search failed: %s", err.Error())
			}

			templateHandler("member/list", message, len(message) == 0, results, w, r)
		})

		//profile
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			m := r.Context().Value(sessionKey).(*model.Member)
			templateHandler("member/profile", "", true, m, w, r)
		})

		r.Post("/profile", func(w http.ResponseWriter, r *http.Request) {
			message := "Update Successful"
			success := true
			m := r.Context().Value(sessionKey).(*model.Member)
			if m.ID > 0 {
				r.ParseForm()

				m.Handle = r.Form.Get("handle")
				m.FirstName = r.Form.Get("firstname")
				m.LastName = r.Form.Get("lastname")
				m.ModifiedBy = m.ID
				m.Enabled = true

				err := m.Save()
				if err != nil {
					log.Error().Err(err).Msg("profile update failed")
					if strings.Contains(err.Error(), "member_handle_key") {
						message = fmt.Sprintf("Update Failed: that handle is already taken")
					} else {
						message = fmt.Sprintf("Update Failed: %s", err.Error())
					}
					success = false
				}
			}

			templateHandler("member/profile", message, success, m, w, r)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			s := chi.URLParam(r, "id")
			id := util.DecodeLink(s)
			me := r.Context().Value(sessionKey).(*model.Member)
			m, err := model.GetMember(id)
			if err != nil {
				log.Warn().Str("id", s).Msg("member not found")
				renderBadRequest(err, w, r)
				return
			}
			if m.ID == 0 {
				renderNotFound(w, r)
				return
			}
			templateHandler("member/details", "", true, MemberData{Follows: me.DoesFollow(id), Results: m}, w, r)
		})

		r.Post("/{id}/follow", func(w http.ResponseWriter, r *http.Request) {
			s := chi.URLParam(r, "id")
			id := util.DecodeLink(s)
			me := r.Context().Value(sessionKey).(*model.Member)
			m, err := model.GetMember(id)
			if err != nil {
				log.Warn().Str("id", s).Msg("member not found")
				renderBadRequest(err, w, r)
				return
			}
			if m.ID == 0 {
				renderNotFound(w, r)
				return
			}

			if !me.DoesFollow(id) && me.ID != id {
				me.Follows = append(me.Follows, m.ID)
				err = me.Save()
				if err != nil {
					log.Error().Err(err).Int64("a", me.ID).Int64("b", m.ID).Msg("a couldn't follow b")
					renderServerError(err, w, r)
					return
				}
				m.Followers = append(m.Followers, me.ID)
				err = m.Save()
				if err != nil {
					log.Error().Err(err).Int64("b", m.ID).Int64("a", me.ID).Msg("b couldn't be followed by a")
					renderServerError(err, w, r)
					return
				}
				log.Info().Int64("a", me.ID).Int64("b", m.ID).Msg("a followed b")
			}
			path := fmt.Sprintf("/member/%s", s)
			http.Redirect(w, r, path, http.StatusSeeOther)
		})

		r.Post("/{id}/unfollow", func(w http.ResponseWriter, r *http.Request) {
			s := chi.URLParam(r, "id")
			id := util.DecodeLink(s)
			me := r.Context().Value(sessionKey).(*model.Member)
			m, err := model.GetMember(id)
			if err != nil {
				log.Warn().Str("id", s).Msg("member not found")
				renderBadRequest(err, w, r)
				return
			}
			if m.ID == 0 {
				renderNotFound(w, r)
				return
			}

			if me.DoesFollow(id) {
				me.Follows = util.Remove(me.Follows, m.ID)
				err = me.Save()
				if err != nil {
					log.Error().Err(err).Int64("a", me.ID).Int64("b", m.ID).Msg("a couldn't unfollow b")
					renderServerError(err, w, r)
					return
				}
				m.Followers = util.Remove(m.Followers, me.ID)
				err = m.Save()
				if err != nil {
					log.Error().Err(err).Int64("b", m.ID).Int64("a", me.ID).Msg("b couldn't be unfollowed by a")
					renderServerError(err, w, r)
					return
				}
				log.Info().Int64("a", me.ID).Int64("b", m.ID).Msg("a unfollowed b")
			}
			path := fmt.Sprintf("/member/%s", s)
			http.Redirect(w, r, path, http.StatusSeeOther)
		})
	})
}

//SessionCtx adds requesting Member to Context
func SessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := getSessionMember(w, r)
		ctx := context.WithValue(r.Context(), sessionKey, &m)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
