package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type MockClient struct {
}

var MockedDoFunction = func(req *http.Request) (*http.Response, error) {
	return nil, errors.New("boom")
}

func (c *MockClient) Do(req *http.Request) (*http.Response, error) {
	return MockedDoFunction(req)
}

func TestBuildUrl(t *testing.T) {
	t.Run("when a word is passed in", func(t *testing.T) {
		want := LinguaApiUrl + "/jejune"
		got := BuildUrl("jejune")

		if want != got {
			t.Errorf("Got %q, want %q", got, want)
		}
	})

	t.Run("when a word is not passed in", func(t *testing.T) {
		want := LinguaApiUrl
		got := BuildUrl("")

		if want != got {
			t.Errorf("Got %q, want %q", got, want)
		}
	})
}

func TestGetDefinitionFromLingua(t *testing.T) {
	t.Run("On a bad request", func(t *testing.T) {
		mockClient := &MockClient{}
		_, err := GetDefinitionFromLingua(mockClient, "jejune")

		if err == nil {
			t.Errorf("Expected to receive an error but did not")
		}
	})

	t.Run("On a happy request", func(t *testing.T) {
		mockClient := &MockClient{}
		body := "my body"
		want := ioutil.NopCloser(bytes.NewReader([]byte(body)))

		MockedDoFunction = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       want,
			}, nil
		}

		got, err := GetDefinitionFromLingua(mockClient, "jejune")

		assertNoError(t, err)

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})
}

func TestBuildDefinitionSummary(t *testing.T) {
	t.Run("It raises an error when the response is empty", func(t *testing.T) {
		linguaResponse := LinguaRobotResponse{}

		_, err := linguaResponse.BuildDefinitionSummary()

		assertError(t, err, NoDefinitionFoundError)
	})

	t.Run("It returns a definition summary", func(t *testing.T) {
		jsonBody := `{
			"entries": [
				{
					"entry": "jejune",
					"lexemes": [
						{
							"partOfSpeech": "adjective",
							"senses": [
								{
									"definition": ""
								}
							]
						}
					],
					"pronunciations": [
						{
							"context": {
								"regions": [
									"United States"
								]
							},
							"transcriptions": [
								{
									"notation": "IPA",
									"transcription": "/jay-june/"
								}
							]
						}
					]
				}
			]
		}`

		response := LinguaRobotResponse{}
		json.Unmarshal([]byte(jsonBody), &response)

		got, err := response.BuildDefinitionSummary()

		assertNoError(t, err)

		want := DefinitionSummary{
			Word:          "jejune",
			Pronunciation: "/jay-june/",
			Definitions: []Definition{
				{
					Meaning:      "",
					PartOfSpeech: "adjective",
				},
			},
		}

		if !reflect.DeepEqual(got.Definitions, want.Definitions) {
			t.Errorf("got %q, want %q", got.Definitions, want.Definitions)
		}
	})
}

func assertNoError(t *testing.T, got error) {
	t.Helper()

	if got != nil {
		t.Fatal("got an error but did not want one")
	}
}

func assertError(t *testing.T, got error, want error) {
	t.Helper()

	if got == nil {
		t.Errorf("Expected an error, but did not get one")
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
