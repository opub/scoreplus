package web

import (
	"context"
	"net/http"

	"github.com/opub/scoreplus/model"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var sessionKey = &contextKey{"Session"}

func routeMembers(r *chi.Mux) {
	r.Route("/member", func(r chi.Router) {
		r.Use(SessionCtx)

		//profile
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			m := r.Context().Value(sessionKey).(*model.Member)
			templateHandler("profile", m, w, r)
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
