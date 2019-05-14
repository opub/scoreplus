package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/opub/scoreplus/model"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/opub/scoreplus/util"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/instagram"
	microsoft "github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/yahoo"
)

//ProviderIndex of supported providers
type ProviderIndex struct {
	Keys []string
	Map  map[string]string
}

var providers = &ProviderIndex{Keys: []string{"amazon", "facebook", "google", "microsoft", "twitter"}, Map: make(map[string]string)}

var providerKey = &contextKey{"Provider"}

func init() {
	for _, p := range providers.Keys {
		registerProvider(p)
	}

	config := util.GetConfig().Session

	store := sessions.NewCookieStore([]byte(config.Secret))
	store.MaxAge(86400 * 30) // 30 days
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = config.Secure
	gothic.Store = store

	gothic.GetProviderName = getProviderName
}

func routeAuth(r *chi.Mux) {
	//login starting point
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		templateHandler("login", "", true, providers, w, r)
	})

	r.Route("/auth", func(r chi.Router) {

		//logout
		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			log.Debug().Msg("routing auth logout")
			gothic.Logout(w, r)
			deleteMemberSession(w)
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
		})

		r.Route("/{provider}", func(r chi.Router) {
			r.Use(ProviderCtx)

			//start authentication
			r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
				// try to get the user without re-authenticating
				log.Debug().Msg("routing auth start")
				if user, err := gothic.CompleteUserAuth(w, r); err == nil {
					findUser(user, w, r)
				} else {
					gothic.BeginAuthHandler(w, r)
				}
			})

			//continue authentication
			r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
				log.Debug().Msg("routing auth callback")
				user, err := gothic.CompleteUserAuth(w, r)
				if err != nil {
					log.Error().Err(err).Msg("user authentication failed")
					render.Render(w, r, ErrServerError(err))
					return
				}
				findUser(user, w, r)
			})
		})
	})
}

func findUser(u goth.User, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nOAUTH_USER:\n%+v\n\n", u)

	m, err := model.GetMemberFromProvider(u.Provider, u.UserID)
	if err != nil {
		log.Error().Err(err).Msg("member lookup failed")
		render.Render(w, r, ErrServerError(err))
		return
	}

	existing := m.ID > 0
	if existing {
		//update last active
		m.LastActive = model.NullTimeNow()
		err := m.Save()
		if err != nil || m.ID == 0 {
			log.Error().Err(err).Msg("member lastactive failed")
			render.Render(w, r, ErrServerError(err))
			return
		}
	} else {
		//first time user that needs to complete registration
		m.FirstName = u.FirstName
		m.LastName = u.LastName
		m.Email = u.Email
		m.Provider = u.Provider
		m.ProviderID = u.UserID
		m.Enabled = false
		err = m.Save()
		if err != nil || m.ID == 0 {
			log.Error().Err(err).Msg("initial member save failed")
			render.Render(w, r, ErrServerError(err))
			return
		}
	}

	log.Info().Str("handle", m.Handle).Bool("existing", existing).Msg("login successful")

	//establish session
	setMemberSession(m, w, r)

	//redirect user
	if m.Enabled {
		w.Header().Set("Location", "/")
	} else {
		w.Header().Set("Location", "/member/profile")
	}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

//ProviderCtx adds requested Provider to Context
func ProviderCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := chi.URLParam(r, "provider")
		ctx := context.WithValue(r.Context(), providerKey, p)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getProviderName(r *http.Request) (string, error) {
	if p, ok := r.Context().Value(providerKey).(string); ok {
		return p, nil
	}
	return "", errors.New("provider not configured")
}

func registerProvider(name string) {
	client := util.GetConfigString("auth." + name + ".key")
	secret := util.GetConfigString("auth." + name + ".secret")
	callback := fmt.Sprintf(util.GetConfigString("auth.callback"), name)

	//TODO Yahoo does not allow custom ports in redirection uri

	providers.Map[name] = strings.Title(name)

	var p goth.Provider
	switch name {
	case "amazon":
		p = amazon.New(client, secret, callback)
	case "facebook":
		p = facebook.New(client, secret, callback, "public_profile", "email")
	case "google":
		p = google.New(client, secret, callback, "https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email")
	case "instagram":
		p = instagram.New(client, secret, callback)
	case "microsoft":
		p = microsoft.New(client, secret, callback)
		p.SetName("microsoft")
	case "twitter":
		p = twitter.New(client, secret, callback)
	case "yahoo":
		p = yahoo.New(client, secret, callback)
	default:
		log.Fatal().Str("name", name).Msg("unsupported provider")
	}
	goth.UseProviders(p)
}
