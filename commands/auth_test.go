package commands

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func drainResponse(resp *http.Response) {
	io.Copy(ioutil.Discard, resp.Body) //nolint:errcheck
	resp.Body.Close()
}

func TestListenCallback(t *testing.T) {
	const csrfState = "5a390de4-b7ed-46b7-bca5-8782eb40302f"
	const want = "1a0da74e-5d29-4f68-9617-3fea5c3cb3db"

	code := listenCallback(context.Background(), 8080, csrfState)

	t.Run("invalid csrf state", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/callback?state=invalid_state&code=%s", want))
		if err != nil {
			t.Fatal(err)
		}
		defer drainResponse(resp)

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("unexpected response status: %d != %d", resp.StatusCode, http.StatusForbidden)
			return
		}

	})
	t.Run("successful completed", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/callback?state=%s&code=%s", csrfState, want))
		if err != nil {
			t.Fatal(err)
		}
		defer drainResponse(resp)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected response status: %d != %d", resp.StatusCode, http.StatusOK)
			return
		}

		got := <-code
		if got != want {
			t.Errorf("unexpected code: %s != %s", got, want)
			return
		}
	})
}
