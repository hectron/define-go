package lingua

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type MockClient struct{}

var MockedDoFunction = func(req *http.Request) (*http.Response, error) {
	return nil, errors.New("boom")
}

func (mc *MockClient) Do(req *http.Request) (*http.Response, error) {
	return MockedDoFunction(req)
}

func TestDefine(t *testing.T) {
	t.Run("When an errors occurs requesting from Lingua", func(t *testing.T) {
		linguaClient := Lingua{HttpClient: &MockClient{}}

		_, err := linguaClient.Define("jejune")

		if err == nil {
			t.Errorf("Expected to receive an error, but did not")
		}
	})

	t.Run("It sends the correct headers", func(t *testing.T) {
		word := "jejune"

		expectedHeadersAndValues := map[string]string{
			"x-rapidapi-key":  "test-key",
			"x-rapidapi-host": ApiHost,
			"x-debug":         word,
		}

		MockedDoFunction = func(req *http.Request) (*http.Response, error) {
			for key, value := range expectedHeadersAndValues {
				got := req.Header.Get(key)

				if got != value {
					t.Errorf("Header \"%s\": got %s, want %s", key, got, value)
				}
			}

			return nil, errors.New("boom")
		}

		linguaClient := Lingua{HttpClient: &MockClient{}, ApiKey: "test-key"}

		linguaClient.Define(word)
	})

	t.Run("It returns a valid summary", func(t *testing.T) {
		MockedDoFunction = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(
					bytes.NewReader(
						[]byte(`{
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
		}`))),
			}, nil
		}

		client := Lingua{HttpClient: &MockClient{}}

		want := Summary{
			Word:          "jejune",
			Pronunciation: "/jay-june/",
			Definitions: []Definition{
				{
					Meaning:      "",
					PartOfSpeech: "adjective",
				},
			},
		}
		got, err := client.Define("jejune")

		if err != nil {
			t.Errorf("Did not expect error, but got %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
