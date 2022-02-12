package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hectron/go-define/lingua"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	AppName = "go-define"
	Version = "1.0.0"
)

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
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
