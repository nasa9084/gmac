package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func TestGetOAuthConfig(t *testing.T) {
	const credFileTmpl = `{
  "installed":{
    "client_id":"%s",
    "project_id":"projectID",
    "auth_uri":"https://example.com/auth",
    "token_uri":"https://example.com/token",
    "auth_provider_x509_cert_url":"https://example.com/oauth2/v1/certs",
    "client_secret":"clientSecret",
    "redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
  }
}`
	setupFs := func() func() {
		fs = afero.NewMemMapFs()
		return func() {
			fs = afero.NewOsFs()
		}
	}

	t.Run("from stdin", func(t *testing.T) {
		teardownFs := setupFs()
		defer teardownFs()

		clientID := "stdin"

		// setup standard input
		var buf bytes.Buffer
		fmt.Fprintf(&buf, credFileTmpl, clientID)
		stdin = &buf
		defer func() { stdin = os.Stdin }()

		// call with "-" means "read from stdin"
		cfg, err := getOAuthConfig("-")
		if err != nil {
			t.Fatal(err)
		}

		if cfg.ClientID != clientID {
			t.Errorf("client id is unexpected value: %s != %s", cfg.ClientID, clientID)
			return
		}

		// check copied credentials.json
		b, err := afero.ReadFile(fs, filepath.Join(configDir, "credentials.json"))
		if err != nil {
			t.Error(err)
		}
		if fmt.Sprintf(credFileTmpl, clientID) != string(b) {
			t.Errorf("credential.json should be copied in config directory, but unexpected contents")
			return
		}
	})

	t.Run("from default path", func(t *testing.T) {
		teardownFs := setupFs()
		defer teardownFs()

		clientID := "default"

		// setup credentials.json in config dir
		fs.MkdirAll(configDir, 0755)
		afero.WriteFile(fs, filepath.Join(configDir, "credentials.json"), []byte(fmt.Sprintf(credFileTmpl, clientID)), 0644)

		// call with empty string means "read from config dir"
		cfg, err := getOAuthConfig("")
		if err != nil {
			t.Fatal(err)
		}

		if cfg.ClientID != clientID {
			t.Errorf("client id is unexpected value: %s != %s", cfg.ClientID, clientID)
			return
		}
	})

	t.Run("from specified path", func(t *testing.T) {
		teardownFs := setupFs()
		defer teardownFs()

		clientID := "tmp"
		credentialsFilepath := filepath.Join("/", "tmp", "credentials.json")

		// setup credentials.json in /tmp
		fs.Mkdir("tmp", 0755)
		afero.WriteFile(fs, credentialsFilepath, []byte(fmt.Sprintf(credFileTmpl, clientID)), 0644)

		cfg, err := getOAuthConfig(credentialsFilepath)
		if err != nil {
			t.Fatal(err)
		}

		if cfg.ClientID != clientID {
			t.Errorf("client id is unexpected value: %s != %s", cfg.ClientID, clientID)
			return
		}

		// check copied credentials.json
		b, err := afero.ReadFile(fs, filepath.Join(configDir, "credentials.json"))
		if err != nil {
			t.Error(err)
		}
		if fmt.Sprintf(credFileTmpl, clientID) != string(b) {
			t.Errorf("credential.json should be copied in config directory, but unexpected contents")
			return
		}
	})
}
