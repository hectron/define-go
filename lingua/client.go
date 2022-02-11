package lingua

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ApiUrl  = "https://lingua-robot.p.rapidapi.com/language/v1/entries/en"
	ApiHost = "lingua-robot.p.rapidapi.com"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Lingua struct {
	HttpClient HTTPClient
	ApiKey     string
}

func (l *Lingua) Define(word string) (Summary, error) {
	var summary Summary

	url := fmt.Sprintf("%s/%s", ApiUrl, word)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return summary, err
	}

	request.Header.Add("x-debug", word)

	for header, value := range l.headers() {
		request.Header.Add(header, value)
	}

	res, err := l.HttpClient.Do(request)

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return summary, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return summary, err
	}

	response := RobotResponse{}
	json.Unmarshal(body, &response)

	summary, err = response.Summary()

	if err != nil {
		return summary, err
	}

	return summary, nil
}

func (l *Lingua) headers() map[string]string {
	return map[string]string{
		"x-rapidapi-host": ApiHost,
		"x-rapidapi-key":  l.ApiKey,
	}
}
