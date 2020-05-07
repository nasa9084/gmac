package commands

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

	"github.com/nasa9084/gmac/log"
)

var oauthScope = []string{
	gmail.GmailLabelsScope,
	gmail.GmailSettingsBasicScope,
}

var fs = afero.NewOsFs()
var stdin io.Reader = os.Stdin

type raw []byte

func (r *raw) UnmarshalYAML(data []byte) error {
	if r == nil {
		return errors.New("UnmarshalYAML on nil pointer")
	}
	*r = append((*r)[0:0], data...)
	return nil
}

func getOAuthConfig(credentialsFilepath string) (*oauth2.Config, error) {
	defaultCredentialsFilepath := filepath.Join(configDir, "credentials.json")

	var r io.Reader
	switch credentialsFilepath {
	case "-": // read from stdin
		log.Vprint("read OAuth config from stdin")
		r = stdin
	case "": // read from default config path
		log.Vprintf("read OAuth config from %s", defaultCredentialsFilepath)
		f, err := fs.Open(defaultCredentialsFilepath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		r = f
	default: // read from specified filepath
		log.Vprintf("read OAuthconfig from %s", credentialsFilepath)
		f, err := fs.Open(credentialsFilepath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		r = f
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	oauthConfig, err := google.ConfigFromJSON(b, oauthScope...)
	if err != nil {
		return nil, err
	}
	if credentialsFilepath != "" {
		log.Vprint("OAuth config is not read from config directory: save into config directory")
		if err := fs.MkdirAll(configDir, 0755); err != nil {
			log.Printf("WARN: cannot create config directory: %s", configDir)
		} else if err := afero.WriteFile(fs, defaultCredentialsFilepath, b, 0644); err != nil {
			log.Printf("WARN: %+v", err)
		}
	}
	return oauthConfig, nil
}

func getToken(refreshToken string) (*oauth2.Token, error) {
	if refreshToken != "" {
		log.Vprint("refresh token is passed")
		return &oauth2.Token{
			RefreshToken: refreshToken,
		}, nil
	}
	tokenFilepath := filepath.Join(configDir, "token.json")
	log.Vprintf("refresh token is not passed, read OAuth token from %s", tokenFilepath)
	b, err := afero.ReadFile(fs, tokenFilepath)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	if err := json.Unmarshal(b, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func open(url string) error {
	var cmd string
	log.Vprintf("detected OS: %s", runtime.GOOS)
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	return exec.Command(cmd, url).Start()
}

func mustWriteString(w io.Writer, s string) {
	if _, err := w.Write([]byte(s)); err != nil {
		panic(err)
	}
}
