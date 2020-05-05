package commands

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

var (
	Version, Revision string
)

var parser = flags.NewParser(&Command{
	ShowVersion: func() error {
		fmt.Printf("Version: %s\nRevision: %s", Version, Revision)
		return &flags.Error{
			Type: flags.ErrHelp,
		}
	},
}, flags.Default)

type Command struct {
	OutputFormat string `short:"o" long:"output" choice:"yaml" choice:"wide"`

	CredentialsFilePath string `short:"c" long:"credentials-file" default:"credentials.json"`
	RefreshToken        string `short:"t" long:"refresh-token"`

	ShowVersion func() error `short:"v" long:"version"`
}

func Run() error {
	if _, err := parser.Parse(); err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			return nil
		}
		return err
	}
	return nil
}

func must(cmd *flags.Command, err error) *flags.Command {
	if err != nil {
		panic(err)
	}
	return cmd
}
