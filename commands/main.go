package commands

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var (
	Version, Revision string
)

var parser = flags.NewParser(&Command{
	ShowVersion: func() error {
		fmt.Printf("Version: %s\nRevision: %s\n", Version, Revision)
		return &flags.Error{
			Type: flags.ErrHelp,
		}
	},
}, flags.HelpFlag)

type Command struct {
	OutputFormat string `short:"o" long:"output" choice:"yaml" choice:"wide"`

	CredentialsFilePath string `short:"c" long:"credentials-file" description:"path to OAuth credentials file"`
	RefreshToken        string `short:"t" long:"refresh-token" description:"OAuth reflesh token"`

	ShowVersion func() error `short:"v" long:"version"`
}

func Run() error {
	if _, err := parser.Parse(); err != nil {
		if fe, ok := err.(*flags.Error); ok {
			switch fe.Type {
			case flags.ErrHelp, flags.ErrCommandRequired:
				if !parser.FindOptionByLongName("version").IsSet() {
					parser.WriteHelp(os.Stdout)
				}
				return nil
			}
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
