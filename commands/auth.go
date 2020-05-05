package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jessevdk/go-flags"
	"golang.org/x/oauth2"
)

var authCommand *flags.Command

func init() {
	authCommand = must(parser.AddCommand("auth", "", "", &AuthCommand{}))
}

type AuthCommand struct {
	CredentialsFilePath string `short:"f" long:"credentials-file" default:"credentials.json"`
	Port                int    `short:"p" long:"port" default:"8080" description:"localhost port to listen callback request"`
}

func (cmd *AuthCommand) Execute(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	oauthConfig, err := getOAuthConfig(cmd.CredentialsFilePath)
	if err != nil {
		return err
	}
	oauthConfig.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", cmd.Port)

	csrfState := uuid.New().String()
	code := listenCallback(ctx, cmd.Port, csrfState)

	if err := open(oauthConfig.AuthCodeURL(csrfState, oauth2.AccessTypeOffline)); err != nil {
		return err
	}

	token, err := oauthConfig.Exchange(ctx, <-code)
	if err != nil {
		return err
	}
	fmt.Printf("Reflesh Token: %s\n", token.RefreshToken)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(token); err != nil {
		return err
	}
	return ioutil.WriteFile("./token.json", buf.Bytes(), 0644)
}

func listenCallback(ctx context.Context, port int, csrfState string) chan string {
	future := make(chan string, 1)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != csrfState {
			log.Print("state mismatch")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Auth Successful Completed, Please Close`)) //nolint:errcheck
		future <- r.FormValue("code")
	})
	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
			log.Printf("error on http.ListenAndServe: %v", err)
		}
	}()
	return future
}
