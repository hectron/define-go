package main

import (
	"bytes"
	"testing"
)

func TestDefinitionSummary(t *testing.T) {
	summary := DefinitionSummary{
		Word:          "jejune",
		Pronunciation: "/jay-june/",
		Definitions: []Definition{
			{
				PartOfSpeech: "adjective",
				Meaning:      "naive, simplistic, and superficial",
				UsageExamples: []string{
					"their entirely predicatable and usually jejune opinions",
					"the poem seems to me rather jejune",
				},
			},
		},
	}

	t.Run("It outputs properly", func(t *testing.T) {
		buffer := bytes.Buffer{}

		summary.Print(&buffer)

		got := buffer.String()
		want := `
jejune (/jay-june/)

[adjective]
  naive, simplistic, and superficial
	e.g. "their entirely predicatable and usually jejune opinions"
	e.g. "the poem seems to me rather jejune"

`

		if got != want {
			t.Errorf("\ngot %q,\nwant %q", got, want)
		}
	})
}
