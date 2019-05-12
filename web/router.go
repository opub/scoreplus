package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

var modelKey = &contextKey{"Model"}

//Start launches web server
func Start() {
	r := chi.NewRouter()

	// our middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(Logger)
	r.Use(Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(render.SetContentType(render.ContentTypeHTML))
	r.Use(middleware.Timeout(60 * time.Second))

	//basic template pages
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templateHandler("home", providers, w, r)
	})
	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		templateHandler("static/home", providers, w, r)
	})
	r.Get("/privacy", func(w http.ResponseWriter, r *http.Request) {
		templateHandler("static/privacy", "", w, r)
	})

	// r.Route("/game", func(r chi.Router) {
	// 	r.Use(render.SetContentType(render.ContentTypeJSON))

	// 	// r.With(paginate).Get("/", listGames) // GET /game
	// 	// r.Post("/", createGame)              // POST /game

	// 	// Subrouters:
	// 	r.Route("/{id}", func(r chi.Router) {
	// 		r.Use(GameCtx)
	// 		r.Get("/", getModel)       // GET /game/123
	// 		r.Delete("/", deleteModel) // DELETE /game/123
	// 		// r.Put("/", updateGame)     // PUT /game/123
	// 	})
	// })

	//additional routing
	routeAuth(r)
	routeMembers(r)
	routeTeams(r)

	//static resources
	config := util.GetConfig()
	wd, _ := os.Getwd()
	filesDir := filepath.Join(wd, config.Path.StaticFiles)
	fileServer(r, "/static", http.Dir(filesDir))
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir(filesDir))
		fs.ServeHTTP(w, r)
	})

	log.Info().Str("version", util.Version).Msg("starting server")
	http.ListenAndServe(":8080", r)
}

func templateHandler(name string, data interface{}, w http.ResponseWriter, r *http.Request) {
	log.Debug().Str("template", name).Msg("template handler")
	if t, ok := Templates[name]; ok {
		err := t.Execute(w, data)
		if err != nil {
			log.Error().Str("template", name).Msg("could not execute template")
			render.Render(w, r, ErrServerError(err))
			return
		}
	} else {
		log.Error().Str("template", name).Msg("template doesn't exist")
	}
}

// fileServer sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// //GameCtx adds requested Game to Context
// func GameCtx(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		s := chi.URLParam(r, "id")
// 		id, err := strconv.ParseInt(s, 10, 64)
// 		if err != nil {
// 			log.Info().Str("id", s).Msg("invalid id")
// 			render.Render(w, r, ErrBadRequest(err))
// 			return
// 		}
// 		game, err := model.GetGame(id)
// 		if err != nil {
// 			log.Warn().Int64("id", id).Msg("game not loaded")
// 			render.Render(w, r, ErrServerError(err))
// 			return
// 		}
// 		if game.ID == 0 {
// 			log.Info().Int64("id", id).Msg("game not found")
// 			render.Render(w, r, ErrNotFound)
// 			return
// 		}
// 		ctx := context.WithValue(r.Context(), modelKey, &game)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// func getModel(w http.ResponseWriter, r *http.Request) {
// 	m := r.Context().Value(modelKey).(model.Model)
// 	render.Render(w, r, NewModelResponse(m))
// }

// func deleteModel(w http.ResponseWriter, r *http.Request) {
// 	m := r.Context().Value(modelKey).(model.Model)
// 	err := m.Delete()
// 	if err != nil {
// 		log.Warn().Msg("delete failed")
// 		render.Render(w, r, ErrServerError(err))
// 		return
// 	}
// 	render.Render(w, r, StatusOK)
// }

// // wrapper methods for

// //ModelResponse result wrapper
// type ModelResponse struct {
// 	Results model.Model `json:"results"`
// }

// //NewModelResponse creates new wrapped ModelResponse
// func NewModelResponse(m model.Model) *ModelResponse {
// 	return &ModelResponse{Results: m}
// }

// //Render interface for JSON rendering
// func (mr *ModelResponse) Render(w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }
