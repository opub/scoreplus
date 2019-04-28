package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
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

	config := util.GetConfig().Auth

	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	store.MaxAge(86400 * 30) // 30 days
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = config.Secure

	gothic.Store = store

	gothic.GetProviderName = getProviderName
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
