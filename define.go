package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hectron/go-define/lingua"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const AppName = "go-define"

// this is injected using ldflags during compile time
// @see Makefile
var Version = ""

func main() {
	app := &cli.App{
		Name:    AppName,
		Version: Version,
		Usage:   "Find the definition of a word",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   filepath.Join(ConfigPath, "config.yml"),
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: func(c *cli.Context) error {
			word := c.Args().Get(0)

			if len(word) == 0 {
				return errors.New("Please provide a word to define")
			}

			LoadConfig(c)

			apiKey := viper.GetString("LINGUA_ROBOT_API_KEY")

			if apiKey == "" {
				return errors.New("Please provide an API key")
			}

			linguaClient := lingua.Lingua{HttpClient: http.DefaultClient, ApiKey: apiKey}

			summary, err := linguaClient.Define(word)

			if err != nil {
				return err
			}

			summary.Print(os.Stdout)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "manage your configuration file",
				Subcommands: []*cli.Command{
					{
						Name:     "init",
						Category: "config",
						Usage:    "Create a new configuration file if doesn't exist",
						Action: func(c *cli.Context) error {
							fileContents := []byte("### Get an API key from https://rapidapi.com/rokish/api/lingua-robot/pricing\nLINGUA_ROBOT_API_KEY: <Key>\n")

							err := os.Mkdir(ConfigPath, 0777)

							if err != nil {
								return err
							}

							err = os.WriteFile(filepath.Join(ConfigPath, "config.yml"), fileContents, 0644)

							if err != nil {
								return err
							}

							fmt.Println("Wrote", filepath.Join(ConfigPath, "config.yml"))

							return nil
						},
					},
					{
						Name:     "edit",
						Category: "config",
						Usage:    "Edit existing configuration",
						Action: func(c *cli.Context) error {
							configFilePath := filepath.Join(ConfigPath, "config.yml")

							editor := os.Getenv("EDITOR")

							if editor == "" {
								return errors.New(fmt.Sprintf("$EDITOR environment variable is not set. Cannot open editor to edit the config file: %s", configFilePath))
							}

							cmd := exec.Command(editor, configFilePath)

							cmd.Stdin = os.Stdin
							cmd.Stdout = os.Stdout

							if err := cmd.Run(); err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:     "rm",
						Category: "config",
						Usage:    "Delete your configuration file",
						Action: func(c *cli.Context) error {
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
