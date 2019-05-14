package web

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/gorilla/securecookie"
	"github.com/opub/scoreplus/model"
	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"
)

const sessionCookie = "spid"

var sc = initSession()

func initSession() *securecookie.SecureCookie {
	config := util.GetConfig().Session
	var h = []byte(config.Hash)
	var b = []byte(config.Block)
	return securecookie.New(h, b)
}

func setMemberSession(m model.Member, w http.ResponseWriter, r *http.Request) {
	expires := time.Now().UTC().AddDate(0, 0, 1)

	value := map[string]string{
		"key": util.RandomString(24),
		"sig": util.EncodeCookie(m.ID, expires),
	}
	encoded, err := sc.Encode(sessionCookie, value)
	if err != nil {
		log.Error().Err(err).Msg("cookie encoding failed")
		render.Render(w, r, ErrServerError(err))
		return
	}

	//set session key in cookie
	cookie := &http.Cookie{
		Name:     sessionCookie,
		Value:    encoded,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

//swallows errors since lack of session is expected and issues with it shouldn't be surfaced to user
func getSessionMember(w http.ResponseWriter, r *http.Request) model.Member {
	m := model.Member{}

	//get session key from cookie
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		log.Warn().Err(err).Msg("couldn't get cookie")
		return m
	}
	value := make(map[string]string)
	err = sc.Decode(sessionCookie, cookie.Value, &value)
	if err != nil {
		log.Warn().Err(err).Msg("couldn't decode cookie")
		return m
	}

	//verify expiration beyond cookie that can be faked
	id, expires := util.DecodeCookie(value["sig"])
	if expires.Before(time.Now()) {
		log.Info().Int64("id", id).Msg("session has expired")
		return m
	}

	//get member
	m, err = model.GetMember(id)
	if err != nil {
		log.Warn().Err(err).Msg("couldn't load member")
	}
	return m
}

func deleteMemberSession(w http.ResponseWriter) {
	//set expired cookie with no data
	cookie := &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
