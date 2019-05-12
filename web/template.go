package web

import (
	"html/template"
	"os"

	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"
)

var baseTemplates = []string{"home", "login"}
var memberTemplates = []string{"list", "details", "profile"}
var teamTemplates = []string{"list", "details", "new"}
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
