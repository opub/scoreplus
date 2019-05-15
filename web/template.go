package web

import (
	"net/http"

	"github.com/opub/scoreplus/model"

	"html/template"
	"os"

	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"
)

type templateGroup struct {
	base  string
	paths []string
}

//Templates that have been loaded into the system
var Templates = make(map[string]*template.Template)

func init() {
	groups := []templateGroup{{"", []string{"home", "login"}}, {"error", []string{"400", "404", "500"}}, {"member", []string{"list", "details", "profile"}},
		{"team", []string{"list", "details", "new"}}, {"venue", []string{"list", "details", "new"}}, {"static", []string{"privacy"}}}

	for _, g := range groups {
		prefix := g.base + "/"
		if len(g.base) == 0 {
			prefix = ""
		}
		for _, p := range g.paths {
			Templates[prefix+p] = parseTemplate(prefix + p)
		}
	}
}

func parseTemplate(name string) *template.Template {
	config := util.GetConfig()
	path := config.Path.Templates
	t, err := template.ParseFiles(path+"layout.html", path+name+".html")
	if err != nil {
		wd, _ := os.Getwd()
		log.Fatal().Str("wd", wd).Str("path", path).Str("name", name).Err(err).Msg("failed to parse template")
		return nil
	}
	return t
}

//TemplateData has data required for template params
type TemplateData struct {
	Message string
	Success bool
	Session bool
	User    model.Member
	Data    interface{}
}

func renderBadRequest(e error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	templateHandler("error/400", "", false, e, w, r)
}

func renderNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templateHandler("error/404", "", false, nil, w, r)
}

func renderServerError(e error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	templateHandler("error/500", "", false, e, w, r)
}

func templateHandler(name string, message string, success bool, data interface{}, w http.ResponseWriter, r *http.Request) {
	log.Debug().Str("template", name).Msg("template handler")
	if t, ok := Templates[name]; ok {
		m := getSessionMember(w, r)
		d := TemplateData{Message: message, Success: success, Session: m.ID != 0, User: m, Data: data}
		err := t.Execute(w, d)
		if err != nil {
			log.Error().Str("template", name).Msg("could not execute template")
			if name != "error/500" {
				renderServerError(err, w, r)
			}
		}
	} else {
		log.Error().Str("template", name).Msg("template doesn't exist")
	}
}
