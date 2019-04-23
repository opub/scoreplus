package web

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/opub/scoreplus/model"
	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type modelContextKey struct{}

var modelKey = modelContextKey{}

//Start launches web server
func Start() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(Logger)
	r.Use(Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	// RESTy routes for "articles" resource
	r.Route("/game", func(r chi.Router) {
		// r.With(paginate).Get("/", listGames) // GET /game
		// r.Post("/", createGame)              // POST /game

		// Subrouters:
		r.Route("/{id}", func(r chi.Router) {
			r.Use(GameCtx)
			r.Get("/", getModel)       // GET /game/123
			r.Delete("/", deleteModel) // DELETE /game/123
			// r.Put("/", updateGame)     // PUT /game/123
		})
	})

	// Mount the admin sub-router
	// r.Mount("/admin", adminRouter())

	log.Info().Str("version", util.Version).Msg("starting server")
	http.ListenAndServe(":8080", r)
}

//GameCtx adds requested Game to Context
func GameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Info().Str("id", s).Msg("invalid id")
			render.Render(w, r, ErrBadRequest(err))
			return
		}
		game, err := model.GetGame(id)
		if err != nil {
			log.Warn().Int64("id", id).Msg("game not loaded")
			render.Render(w, r, ErrServerError(err))
			return
		}
		if game.ID == 0 {
			log.Info().Int64("id", id).Msg("game not found")
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), modelKey, game)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getModel(w http.ResponseWriter, r *http.Request) {
	m := r.Context().Value(modelKey).(*model.Base)
	render.Render(w, r, m)
}

func deleteModel(w http.ResponseWriter, r *http.Request) {
	m := r.Context().Value(modelKey).(*model.Base)
	err := m.Delete()
	if err != nil {
		log.Warn().Int64("id", m.ID).Msg("delete failed")
		render.Render(w, r, ErrServerError(err))
		return
	}
	render.Render(w, r, StatusOK)
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}
