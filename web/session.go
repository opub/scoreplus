package web

import (
	"net/http"
	"strconv"
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
	key := util.RandomString(32)
	expires := time.Now().AddDate(0, 0, 1) // 1 day

	value := map[string]string{
		"key": key,
		"exp": strconv.FormatInt(expires.UTC().Unix(), 10),
		"usr": strconv.FormatInt(m.ID, 10),
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
	expires, err := strconv.ParseInt(value["exp"], 10, 64)
	if err != nil {
		return m, err
	}
	if time.Unix(expires, 0).After(time.Now()) {
		return m, errors.New("session has expired")
	}

	//get member
	id, err := strconv.ParseInt(value["usr"], 10, 64)
	if err != nil {
		return m, err
	}
	return model.GetMember(id)
}
