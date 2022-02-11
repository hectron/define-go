package lingua

import (
	"bytes"
	"testing"
)

func TestSummary(t *testing.T) {
	testCases := []struct {
		description, want string
		input             Summary
	}{
		{
			description: "It outputs a proper definition",
			input: Summary{
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
			},
			want: `
jejune (/jay-june/)

[adjective]
  naive, simplistic, and superficial
	e.g. "their entirely predicatable and usually jejune opinions"
	e.g. "the poem seems to me rather jejune"

`,
		},
		{
			description: "It does not include a pronunciation in the output if it doesn't exist",
			input: Summary{
				Word:          "jejune",
				Pronunciation: "",
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
			},
			want: `
jejune

[adjective]
  naive, simplistic, and superficial
	e.g. "their entirely predicatable and usually jejune opinions"
	e.g. "the poem seems to me rather jejune"

`,
		},
		{
			description: "It includes multiple definitions",
			input: Summary{
				Word:          "jejune",
				Pronunciation: "",
				Definitions: []Definition{
					{
						PartOfSpeech: "adjective",
						Meaning:      "naive, simplistic, and superficial",
						UsageExamples: []string{
							"their entirely predicatable and usually jejune opinions",
						},
					},
					{
						PartOfSpeech: "adjective",
						Meaning:      "(of ideas or writings) dry and uninteresting.",
						UsageExamples: []string{
							"the poem seems to me rather jejune",
						},
					},
				},
			},
			want: `
jejune

[adjective]
  naive, simplistic, and superficial
	e.g. "their entirely predicatable and usually jejune opinions"

[adjective]
  (of ideas or writings) dry and uninteresting.
	e.g. "the poem seems to me rather jejune"

`,
		},
		{
			description: "It works without usage examples",
			input: Summary{
				Word:          "jejune",
				Pronunciation: "",
				Definitions: []Definition{
					{
						PartOfSpeech: "adjective",
						Meaning:      "naive, simplistic, and superficial",
					},
					{
						PartOfSpeech: "adjective",
						Meaning:      "(of ideas or writings) dry and uninteresting.",
					},
				},
			},
			want: `
jejune

[adjective]
  naive, simplistic, and superficial

[adjective]
  (of ideas or writings) dry and uninteresting.

`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			buffer := bytes.Buffer{}
			tc.input.Print(&buffer)
			got := buffer.String()

			if got != tc.want {
				t.Errorf("\ngot  %q,\nwant %q", got, tc.want)
			}
		})
	}
}
