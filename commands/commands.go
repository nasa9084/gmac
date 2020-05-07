package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"

	"github.com/nasa9084/gmac/log"
)

var (
	Version, Revision string
)

var parser = flags.NewParser(&Command{
	Verbose: func() error {
		log.SetVerbose(true)
		return nil
	},

	ShowVersion: func() error {
		fmt.Printf("Version: %s\nRevision: %s\n", Version, Revision)
		return &flags.Error{
			Type: flags.ErrHelp,
		}
	},
}, flags.HelpFlag)

var configDir string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Print("WARN: cannot detect user home directory")
	}
	configDir = filepath.Join(home, ".gmac")
}

type Command struct {
	OutputFormat string `short:"o" long:"output" choice:"yaml" choice:"wide"`

	CredentialsFilePath string `short:"c" long:"credentials-file" description:"path to OAuth credentials file"`
	RefreshToken        string `short:"t" long:"refresh-token" description:"OAuth reflesh token"`

	Verbose func() error `long:"verbose" description:"show verbose log"`

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
