package commands

import (
	"context"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/nasa9084/gmac/encoder"
	"github.com/nasa9084/gmac/gmail"
)

var (
	getCommand       *flags.Command
	getFilterCommand *flags.Command
)

func init() {
	getCommand = must(parser.AddCommand("get", "Get resources", "Get resources", &GetCommand{}))
	getFilterCommand = must(getCommand.AddCommand("filter", "", "", &GetFilterCommand{}))
	getFilterCommand.Aliases = []string{"filters"}
}

type GetCommand struct {
}

type GetFilterCommand struct {
}

func (cmd *GetFilterCommand) Execute([]string) error {
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
	filters, err := c.ListFilters(ctx)

	if err := encoder.NewFilterEncoder(os.Stdout, cmd.OutputFormat()).Encode(filters); err != nil {
		return err
	}

	return nil
}

func (*GetFilterCommand) CredentialsFilePath() string {
	val := getFilterCommand.FindOptionByLongName("credentials-file").Value()
	if val == nil {
		return ""
	}
	return val.(string)
}

func (*GetFilterCommand) RefreshToken() string {
	val := getFilterCommand.FindOptionByLongName("refresh-token").Value()
	if val == nil {
		return ""
	}
	return val.(string)
}

func (*GetFilterCommand) OutputFormat() string {
	val := getFilterCommand.FindOptionByLongName("output").Value()
	if val == nil {
		return ""
	}
	return val.(string)
}
