package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var NoDefinitionFoundError = errors.New("No definition found")

const (
	LinguaApiUrl  = "https://lingua-robot.p.rapidapi.com/language/v1/entries/en"
	LinguaApiHost = "lingua-robot.p.rapidapi.com"
)

var (
	LinguaRobotApiKey string = os.Getenv("LINGUA_ROBOT_API_KEY")
	LinguaApiHeaders         = map[string]string{
		"x-rapidapi-host": LinguaApiHost,
		"x-rapidapi-key":  LinguaRobotApiKey,
	}
)

// I cheated because this was auto-generated using https://mholt.github.io/json-to-go/
// This is kind-of insane
type LinguaRobotResponse struct {
	Entries []struct {
		Entry   string `json:"entry"`
		Lexemes []struct {
			PartOfSpeech string `json:"partOfSpeech"`
			Senses       []struct {
				Definition    string   `json:"definition"`
				UsageExamples []string `json:"usageExamples,omitempty"`
				Context       struct {
					Regions []string `json:"regions"`
				} `json:"context,omitempty"`
			} `json:"senses"`
		} `json:"lexemes"`
		Pronunciations []struct {
			Context struct {
				Regions []string `json:"regions"`
			} `json:"context"`
			Transcriptions []struct {
				Notation      string `json:"notation"`
				Transcription string `json:"transcription"`
			} `json:"transcriptions"`
		} `json:"pronunciations"`
	} `json:"entries"`
}

func LookUp(client HTTPClient, word string) (LinguaRobotResponse, error) {
	response, err := GetDefinitionFromLingua(client, word)

	if err != nil {
		return LinguaRobotResponse{}, errors.New("Unable to retrieve definition from Lingua")
	}

	body, _ := ioutil.ReadAll(response)

	linguaResponse := LinguaRobotResponse{}
	json.Unmarshal([]byte(body), &linguaResponse)

	return linguaResponse, nil
}

func BuildUrl(word string) string {
	if word == "" {
		return LinguaApiUrl
	}

	return fmt.Sprintf("%s/%s", LinguaApiUrl, word)
}

func GetDefinitionFromLingua(client HTTPClient, word string) (io.ReadCloser, error) {
	url := BuildUrl(word)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	for header, value := range LinguaApiHeaders {
		request.Header.Add(header, value)
	}

	res, err := client.Do(request)

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// This is the core of the tool
//
// Takes a JSON response, and converts it into an DefinitionSummary
//
// This makes a few assumptions, such as only wanting pronunciation from the U.S.,
// and only checking the first entry that is returned from Lingua
func (response LinguaRobotResponse) BuildDefinitionSummary() (DefinitionSummary, error) {
	if len(response.Entries) == 0 {
		return DefinitionSummary{}, NoDefinitionFoundError
	}

	entry := response.Entries[0]
	word := entry.Entry
	var transcription string
	var definitions []Definition

	for i := 0; i < len(entry.Pronunciations); i++ {
		pronunciation := entry.Pronunciations[i]

		for j := 0; j < len(pronunciation.Context.Regions); j++ {
			region := pronunciation.Context.Regions[j]

			if region == "United States" {
				transcription = pronunciation.Transcriptions[0].Transcription
			}
		}
	}

	for i := 0; i < len(entry.Lexemes); i++ {
		lex := entry.Lexemes[i]

		hasExamples := false

		for j := 0; j < len(lex.Senses); j++ {
			sense := lex.Senses[j]

			if len(sense.UsageExamples) > 0 {
				hasExamples = true

				definitions = append(definitions, Definition{
					lex.PartOfSpeech,
					sense.Definition,
					sense.UsageExamples[0:2],
				})
			}
		}

		// there might be a chance that we don't have a definition
		if !hasExamples {
			sense := lex.Senses[0]

			definitions = append(definitions, Definition{
				lex.PartOfSpeech,
				sense.Definition,
				sense.UsageExamples,
			})
		}
	}

	return DefinitionSummary{word, transcription, definitions}, nil
}

func Define(word string) (DefinitionSummary, error) {
	client := http.DefaultClient
	response, err := LookUp(client, word)

	if err != nil {
		return DefinitionSummary{}, err
	}

	return response.BuildDefinitionSummary()
}

func main() {
	app := &cli.App{
		Name:  "define",
		Usage: "Find the definition of a word",
		Action: func(c *cli.Context) error {
			word := c.Args().Get(0)

			summary, err := Define(word)

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
