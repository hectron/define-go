package main

import (
	"fmt"
	"io"
)

type Definition struct {
	PartOfSpeech, Meaning string
	UsageExamples         []string
}

type DefinitionSummary struct {
	Word, Pronunciation string
	Definitions         []Definition
}

func (d DefinitionSummary) Print(writer io.Writer) {
	fmt.Fprintf(writer, "\n")

	if len(d.Pronunciation) > 0 {
		fmt.Fprintf(writer, "%s (%s)\n", d.Word, d.Pronunciation)
	} else {
		fmt.Fprintf(writer, "%s\n", d.Word)
	}

	for i := 0; i < len(d.Definitions); i++ {
		definition := d.Definitions[i]

		fmt.Fprintf(writer, "\n")
		fmt.Fprintf(writer, "[%s]\n  %s\n", definition.PartOfSpeech, definition.Meaning)

		for j := 0; j < len(definition.UsageExamples); j++ {
			example := definition.UsageExamples[j]

			if len(example) > 0 {
				fmt.Fprintf(writer, "\te.g. \"%s\"\n", example)
			}
		}
	}

	fmt.Fprintf(writer, "\n")
}
