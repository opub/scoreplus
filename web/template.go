package web

import (
	"net/http"

	"github.com/opub/scoreplus/model"

	"html/template"
	"os"

	"github.com/go-chi/render"
	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"
)

var baseTemplates = []string{"home", "login"}
var memberTemplates = []string{"list", "details", "profile"}
var teamTemplates = []string{"list", "details", "new"}
var venueTemplates = []string{"list", "details", "new"}
var staticTemplates = []string{"privacy"}

//Templates that have been loaded into the system
var Templates = make(map[string]*template.Template)

func init() {
	for _, n := range baseTemplates {
		Templates[n] = parseTemplate(n)
	}
	for _, n := range memberTemplates {
		Templates["member/"+n] = parseTemplate("member/" + n)
	}
	for _, n := range teamTemplates {
		Templates["team/"+n] = parseTemplate("team/" + n)
	}
	for _, n := range venueTemplates {
		Templates["venue/"+n] = parseTemplate("venue/" + n)
	}
	for _, n := range staticTemplates {
		Templates["static/"+n] = parseTemplate("static/" + n)
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

func templateHandler(name string, message string, success bool, data interface{}, w http.ResponseWriter, r *http.Request) {
	log.Debug().Str("template", name).Msg("template handler")
	if t, ok := Templates[name]; ok {
		m := getSessionMember(w, r)
		d := TemplateData{Message: message, Success: success, Session: m.ID != 0, User: m, Data: data}
		err := t.Execute(w, d)
		if err != nil {
			log.Error().Str("template", name).Msg("could not execute template")
			render.Render(w, r, ErrServerError(err))
			return
		}
	} else {
		log.Error().Str("template", name).Msg("template doesn't exist")
	}
}
