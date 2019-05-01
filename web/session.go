package web

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/gorilla/securecookie"
	"github.com/opub/scoreplus/model"
	"github.com/opub/scoreplus/util"
	"github.com/pkg/errors"
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
		"sig": util.EncodeID(m.ID, expires),
	}
	encoded, err := sc.Encode(sessionCookie, value)
	if err != nil {
		log.Error().Err(err).Msg("cookie encoding failed")
		render.Render(w, r, ErrServerError(err))
		return
	}

	//set session key in cookie
	cookie := &http.Cookie{
		Name:    sessionCookie,
		Value:   encoded,
		Expires: expires,
		Path:    "/",
	}
	http.SetCookie(w, cookie)
}

func getSessionMember(w http.ResponseWriter, r *http.Request) (model.Member, error) {
	m := model.Member{}

	//get session key from cookie
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return m, err
	}
	value := make(map[string]string)
	err = sc.Decode(sessionCookie, cookie.Value, &value)
	if err != nil {
		return m, err
	}

	//verify expiration beyond cookie that can be faked
	id, expires := util.DecodeID(value["sig"])
	if expires.Before(time.Now()) {
		return m, errors.New("session has expired")
	}

	//get member
	return model.GetMember(id)
}
