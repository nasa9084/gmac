package commands

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type raw []byte

func (r *raw) UnmarshalYAML(data []byte) error {
	if r == nil {
		return errors.New("UnmarshalYAML on nil pointer")
	}
	*r = append((*r)[0:0], data...)
	return nil
}

func openOrStdin(filename string) (io.ReadCloser, error) {
	if filename == "-" {
		return ioutil.NopCloser(os.Stdin), nil
	}
	return os.Open(filename)
}

func getOAuthConfig(credentialsFilePath string) (*oauth2.Config, error) {
	rc, err := openOrStdin(credentialsFilePath)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	oauthConfig, err := google.ConfigFromJSON(b, gmail.GmailLabelsScope, gmail.GmailSettingsBasicScope)
	if err != nil {
		return nil, err
	}
	return oauthConfig, nil
}

func getToken(refreshToken string) (*oauth2.Token, error) {
	if refreshToken != "" {
		return &oauth2.Token{
			RefreshToken: refreshToken,
		}, nil
	}
	b, err := ioutil.ReadFile("./token.json")
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
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	return exec.Command(cmd, url).Start()
}
