package lingua

import (
	"errors"
	"net/http"
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
}
