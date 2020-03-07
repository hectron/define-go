package main

import (
  "errors"
  "fmt"
  "os"

  "github.com/urfave/cli/v2"
)

func runCLI(args []string) {
  app := &cli.App{
    Name: "define",
    Version: "v1.0.0",
    Authors: []*cli.Author{
      &cli.Author{
        Name: "Hector Rios",
        Email: "that.hector@gmail.com",
      },
    },
    Usage: "Find the definition of a word",
    UsageText: "define [WORD]",
    Action: func(c *cli.Context) error {
      if c.NArg() == 0 {
        cli.ShowAppHelp(c)
        return errors.New("No argument provided")
      }

      word := c.Args().Get(0)

      for i := 1; i < c.Args().Len(); i++ {
        word += "+"
        word += c.Args().Get(i)
      }

      fmt.Println(word)
      //Define(word)

      return nil
    },
  }

  err := app.Run(os.Args)

  if err != nil {
    os.Exit(1)
  }
}

func main() {
  runCLI(os.Args)
}
