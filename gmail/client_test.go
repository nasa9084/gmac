package gmail

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

func testOAuth(t *testing.T) (*httptest.Server, *oauth2.Config, *oauth2.Token) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("access to unknown path of OAuth server: %s", r.URL.Path)
		}
	}))
	cfg := oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: srv.URL + "/token",
		},
	}
	return srv, &cfg, &oauth2.Token{AccessToken: "access-token"}
}
