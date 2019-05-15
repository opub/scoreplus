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
		templateHandler("home", "", true, providers, w, r)
	})
	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		templateHandler("static/home", "", true, providers, w, r)
	})
	r.Get("/privacy", func(w http.ResponseWriter, r *http.Request) {
		templateHandler("static/privacy", "", true, nil, w, r)
	})

	//additional routing
	routeAuth(r)
	routeMembers(r)
	routeTeams(r)
	routeVenues(r)

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
