package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/opub/scoreplus/model"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var sessionKey = &contextKey{"Session"}

//ActionResult has message to return about action that just happened
type ActionResult struct {
	Message string
	Success bool
	Data    interface{}
}

func routeMembers(r *chi.Mux) {
	r.Route("/member", func(r chi.Router) {
		r.Use(SessionCtx)

		//profile
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			m := r.Context().Value(sessionKey).(*model.Member)
			templateHandler("profile", ActionResult{Data: m}, w, r)
		})

		r.Post("/profile", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			m := r.Context().Value(sessionKey).(*model.Member)

			m.Handle = r.Form.Get("handle")
			m.FirstName = r.Form.Get("firstname")
			m.LastName = r.Form.Get("lastname")
			m.ModifiedBy = m.ID
			m.Enabled = true

			a := ActionResult{Message: "Update successful!", Success: true, Data: m}

			err := m.Save()
			if err != nil {
				log.Warn().Err(err).Msg("profile update failed")
				if strings.Contains(err.Error(), "member_handle_key") {
					a.Message = fmt.Sprintf("Update failed: that handle is already taken")
				} else {
					a.Message = fmt.Sprintf("Update failed: %s", err.Error())
				}
				a.Success = false
			}

			templateHandler("profile", a, w, r)
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
