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
	"github.com/go-chi/render"
)

var sessionKey = &contextKey{"Session"}

//MemberData has data required for template params
type MemberData struct {
	Message string
	Success bool
	Session bool
	Follows bool
	Data    interface{}
}

func routeMembers(r *chi.Mux) {
	r.Route("/member", func(r chi.Router) {
		r.Use(SessionCtx)

		//list
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			templateHandler("member/list", MemberData{}, w, r)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			a := MemberData{Message: "Search successful!", Success: true}

			search := r.Form.Get("search")
			results, err := model.SearchMembers(search)
			if err != nil {
				log.Error().Err(err).Msg("member search failed")
				a.Message = fmt.Sprintf("Search failed: %s", err.Error())
				a.Success = false
			} else {
				a.Data = results
			}

			templateHandler("member/list", a, w, r)
		})

		//profile
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			m := r.Context().Value(sessionKey).(*model.Member)
			templateHandler("member/profile", MemberData{Data: m}, w, r)
		})

		r.Post("/profile", func(w http.ResponseWriter, r *http.Request) {
			m := r.Context().Value(sessionKey).(*model.Member)
			a := MemberData{Message: "Update successful!", Success: true, Data: m}

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
						a.Message = fmt.Sprintf("Update failed: that handle is already taken")
					} else {
						a.Message = fmt.Sprintf("Update failed: %s", err.Error())
					}
					a.Success = false
				}
			}

			templateHandler("member/profile", a, w, r)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			s := chi.URLParam(r, "id")
			id := util.DecodeLink(s)
			me := r.Context().Value(sessionKey).(*model.Member)
			a := MemberData{Follows: me.DoesFollow(id), Session: (me.ID != 0 && me.ID != id)}
			m, err := model.GetMember(id)
			if err != nil {
				log.Warn().Str("id", s).Msg("member not found")
				render.Render(w, r, ErrBadRequest(err))
				return
			}
			if m.ID == 0 {
				render.Render(w, r, ErrNotFound)
				return
			}
			a.Data = m
			templateHandler("member/details", a, w, r)
		})

		r.Post("/{id}/follow", func(w http.ResponseWriter, r *http.Request) {
			s := chi.URLParam(r, "id")
			id := util.DecodeLink(s)
			me := r.Context().Value(sessionKey).(*model.Member)
			m, err := model.GetMember(id)
			if err != nil {
				log.Warn().Str("id", s).Msg("member not found")
				render.Render(w, r, ErrBadRequest(err))
				return
			}
			if m.ID == 0 {
				render.Render(w, r, ErrNotFound)
				return
			}

			if !me.DoesFollow(id) && me.ID != id {
				me.Follows = append(me.Follows, m.ID)
				err = me.Save()
				if err != nil {
					log.Error().Err(err).Int64("a", me.ID).Int64("b", m.ID).Msg("a couldn't follow b")
					render.Render(w, r, ErrServerError(err))
					return
				}
				m.Followers = append(m.Followers, me.ID)
				err = m.Save()
				if err != nil {
					log.Error().Err(err).Int64("b", m.ID).Int64("a", me.ID).Msg("b couldn't be followed by a")
					render.Render(w, r, ErrServerError(err))
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
				render.Render(w, r, ErrBadRequest(err))
				return
			}
			if m.ID == 0 {
				render.Render(w, r, ErrNotFound)
				return
			}

			if me.DoesFollow(id) {
				me.Follows = util.Remove(me.Follows, m.ID)
				err = me.Save()
				if err != nil {
					log.Error().Err(err).Int64("a", me.ID).Int64("b", m.ID).Msg("a couldn't unfollow b")
					render.Render(w, r, ErrServerError(err))
					return
				}
				m.Followers = util.Remove(m.Followers, me.ID)
				err = m.Save()
				if err != nil {
					log.Error().Err(err).Int64("b", m.ID).Int64("a", me.ID).Msg("b couldn't be unfollowed by a")
					render.Render(w, r, ErrServerError(err))
					return
				}
				log.Info().Int64("a", me.ID).Int64("b", m.ID).Msg("a unfollowed b")
			}
			path := fmt.Sprintf("/member/%s", s)
			http.Redirect(w, r, path, http.StatusSeeOther)
		})
	})
}

//SessionCtx adds requested Member to Context
func SessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m, err := getSessionMember(w, r)
		if err != nil {
			log.Warn().Err(err).Stack().Msg("session not loaded")
			render.Render(w, r, ErrUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), sessionKey, &m)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
