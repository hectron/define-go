package lingua

import "errors"

var NoDefinitionFoundError = errors.New("No definition found")

type RobotResponse struct {
	Entries []ResponseEntry `json:"entries"`
}

type ResponseEntry struct {
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
}

func (rr *RobotResponse) Summary() (Summary, error) {
	if len(rr.Entries) == 0 {
		return Summary{}, NoDefinitionFoundError
	}

	entry := rr.Entries[0]

	return Summary{
		entry.Entry,
		entry.FindTranscription(),
		entry.FindDefinitions(),
	}, nil
}

func (re *ResponseEntry) FindTranscription() string {
	for _, pronunciation := range re.Pronunciations {
		hasTranscriptions := len(pronunciation.Transcriptions) > 0

		if !hasTranscriptions {
			return ""
		}

		for _, region := range pronunciation.Context.Regions {
			if region == "United States" {
				return pronunciation.Transcriptions[0].Transcription
			}
		}
	}

	return ""
}

func (re *ResponseEntry) FindDefinitions() (definitions []Definition) {
	// prefer definitions that have examples
	// if we still don't have any, then use the first definition
	for _, lex := range re.Lexemes {
		hasExamples := false
		for _, sense := range lex.Senses {
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

	return definitions
}
