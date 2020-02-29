package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "io/ioutil"
  "log"
  "os"

  "github.com/urfave/cli/v2"
)

const LinguaRobotApiKey = os.Getenv("LINGUA_ROBOT_API_KEY")

// I cheated because this was auto-generated using https://mholt.github.io/json-to-go/
// This is kind-of insane
type LinguaRobotResponse struct {
  Entries []struct {
    Entry           string     `json:"entry"`
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
      Audio struct {
        URL       string `json:"url"`
      } `json:"audio,omitempty"`
    } `json:"pronunciations"`
  } `json:"entries"`
}

type Definition struct {
  PartOfSpeech, Meaning string
  UsageExamples []string
}

type Output struct {
  Word, Pronunciation, PronunciationUrl string
  Definitions []Definition
}

func (o Output) Print() {
  fmt.Println("")
  if len(o.Pronunciation) > 0 {
    fmt.Printf("%s (%s)\n", o.Word, o.Pronunciation)
  } else {
    fmt.Printf("%s\n", o.Word)
  }

  for i := 0; i < len(o.Definitions); i++ {
    definition := o.Definitions[i]

    fmt.Println("")
    fmt.Printf("[%s]\n  %s\n", definition.PartOfSpeech, definition.Meaning)

    for j := 0; j < len(definition.UsageExamples); j++ {
      example := definition.UsageExamples[j]

      if len(example) > 0 {
        fmt.Printf("\te.g. \"%s\"\n", example)
      }
    }
  }

  fmt.Println("")
}

func lookUp(word string) LinguaRobotResponse {
  apiUrl := "https://lingua-robot.p.rapidapi.com/language/v1/entries/en"
  hostForHeader := "lingua-robot.p.rapidapi.com"
  url := fmt.Sprintf("%s/%s", apiUrl, word)

  request, _ := http.NewRequest("GET", url, nil)
  request.Header.Add("x-rapidapi-host", hostForHeader)
  request.Header.Add("x-rapidapi-key", LinguaRobotApiKey)

  res, _ := http.DefaultClient.Do(request)

  defer res.Body.Close()
  body, _ := ioutil.ReadAll(res.Body)

  response := LinguaRobotResponse{}
  json.Unmarshal([]byte(body), &response)

  return response
}


// This is the core of the tool
//
// Takes a JSON response, and converts it into an Output
//
// This makes a few assumptions, such as only wanting pronunciation from the U.S.,
// and only checking the first entry that is returned from Lingua
func (response LinguaRobotResponse) BuildOutput() Output {
  entry := response.Entries[0]
  word := entry.Entry
  var transcription string
  var audioUrl string
  var definitions []Definition

  for i := 0; i < len(entry.Pronunciations); i++ {
    pronunciation := entry.Pronunciations[i]

    for j := 0; j < len(pronunciation.Context.Regions); j++ {
      region := pronunciation.Context.Regions[j]

      if region == "United States" {
        audioUrl = pronunciation.Audio.URL
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
    if hasExamples == false {
      sense := lex.Senses[0]

      definitions = append(definitions, Definition{
        lex.PartOfSpeech,
        sense.Definition,
        sense.UsageExamples,
      })
    }
  }

  return Output{word, transcription, audioUrl, definitions}
}


func Define(word string) {
  response := lookUp(word)
  output := response.BuildOutput()
  output.Print()
}

func main() {
  app := &cli.App{
    Name: "define",
    Usage: "Find the definition of a word",
    Action: func(c *cli.Context) error {
      var word string

      word = c.Args().Get(0)

      Define(word)

      return nil
    },
  }

  err := app.Run(os.Args)

  if err != nil {
    log.Fatal(err)
  }
}
