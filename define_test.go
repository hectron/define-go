package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
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

		if err != nil {
			t.Errorf("Expected to receive an error but did not")
		}

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})
}
