package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/jessevdk/go-flags"

	"github.com/nasa9084/gmac/gmail"
	"github.com/nasa9084/gmac/log"
)

var applyCommand *flags.Command

func init() {
	applyCommand = must(parser.AddCommand("apply", "Apply resource", "Apply resource", &ApplyCommand{}))
}

type ApplyCommand struct {
	Target                string `short:"f" long:"filename" required:"yes"`
	ApplyToExistingEmails bool   `short:"e" long:"apply-to-existing"`
}

func (cmd *ApplyCommand) Execute([]string) error {
	var r io.Reader
	switch cmd.Target {
	case "-":
		r = os.Stdin
	default:
		f, err := os.Open(cmd.Target)
		if err != nil {
			return err
		}
		defer f.Close()

		r = f
	}

	var proxy struct {
		Kind string         `yaml:"kind"`
		Rest map[string]raw `yaml:",inline"`
	}
	log.Println("unmarshalYAML")
	if err := yaml.NewDecoder(r).Decode(&proxy); err != nil {
		return err
	}
	if proxy.Kind == "" {
		return errors.New("kind is not found")
	}

	switch proxy.Kind {
	case gmail.ResourceTypeFilter:
		return cmd.applyFilter(proxy.Rest["filters"])
	}

	return fmt.Errorf("unknown resource kind: %s", proxy.Kind)
}

func (cmd *ApplyCommand) applyFilter(data []byte) error {
	if len(data) == 0 {
		return errors.New("required key `filters` not found")
	}
	var filters []gmail.Filter
	if err := yaml.Unmarshal(data, &filters); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	oauthConfig, err := getOAuthConfig(cmd.CredentialsFilePath())
	if err != nil {
		return err
	}

	token, err := getToken(cmd.RefreshToken())
	if err != nil {
		return err
	}

	c, err := gmail.New(ctx, oauthConfig, token)
	if err != nil {
		return err
	}

	log.Printf("Delete all filters...")
	if err := c.DeleteAllFilter(ctx); err != nil {
		return err
	}
	for _, filter := range filters {
		log.Printf("Create filter: %s", filter.String())
		if err := c.CreateFilter(ctx, filter); err != nil {
			return err
		}

		if cmd.ApplyToExistingEmails {
			log.Printf("Apply filter %s to existing emails", filter.String())
			if err := c.ApplyLabelToExistingEmail(ctx, filter); err != nil {
				return err
			}
		}
	}

	return nil
}

func (*ApplyCommand) CredentialsFilePath() string {
	val := applyCommand.FindOptionByLongName("credentials-file").Value()
	if val == nil {
		return ""
	}
	return val.(string)
}

func (*ApplyCommand) RefreshToken() string {
	val := applyCommand.FindOptionByLongName("refresh-token").Value()
	if val == nil {
		return ""
	}
	return val.(string)
}
