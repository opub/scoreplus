package web

import (
	"context"
	"net/http"
	"os"
	"sort"

	"github.com/go-chi/chi"
	"github.com/markbates/goth/gothic"
	"github.com/pkg/errors"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/yahoo"
)

//ProviderIndex of supported providers
type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

var providerIndex *ProviderIndex

var providerKey = &contextKey{"Provider"}

func init() {
	goth.UseProviders(
		//TODO set keys secrets and URLs
		amazon.New(os.Getenv("AMAZON_KEY"), os.Getenv("AMAZON_SECRET"), "http://localhost:3000/auth/amazon/callback"),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
		instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), "http://localhost:3000/auth/instagram/callback"),
		microsoftonline.New(os.Getenv("MICROSOFT_KEY"), os.Getenv("MICROSOFT_SECRET"), "http://localhost:3000/auth/microsoft/callback"),
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),

		//Pointed localhost.com to http://localhost:3000/auth/yahoo/callback through proxy as yahoo
		// does not allow to put custom ports in redirection uri
		yahoo.New(os.Getenv("YAHOO_KEY"), os.Getenv("YAHOO_SECRET"), "http://localhost.com"),
	)

	m := make(map[string]string)
	m["amazon"] = "Amazon"
	m["facebook"] = "Facebook"
	m["google"] = "Google"
	m["instagram"] = "Instagram"
	m["microsoft"] = "Microsoft"
	m["twitter"] = "Twitter"
	m["yahoo"] = "Yahoo"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	providerIndex = &ProviderIndex{Providers: keys, ProvidersMap: m}

	gothic.GetProviderName = getProviderName
}

//ProviderCtx adds requested Provider to Context
func ProviderCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := chi.URLParam(r, "provider")
		ctx := context.WithValue(r.Context(), providerKey, &p)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getProviderName(r *http.Request) (string, error) {
	if p, ok := r.Context().Value(providerKey).(string); ok {
		return p, nil
	}
	return "", errors.New("you must select a provider")
}
