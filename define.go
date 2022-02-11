package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/hectron/go-define/lingua"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "define",
		Usage: "Find the definition of a word",
		Action: func(c *cli.Context) error {
			word := c.Args().Get(0)

			if len(word) == 0 {
				return errors.New("Please provide a word to define")
			}

			linguaClient := lingua.Lingua{
				HttpClient: http.DefaultClient,
				ApiKey:     os.Getenv("LINGUA_ROBOT_API_KEY"),
			}

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
