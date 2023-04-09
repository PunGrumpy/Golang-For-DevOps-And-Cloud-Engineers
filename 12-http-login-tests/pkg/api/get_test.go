package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockClient struct {
	GetResponseOutput  *http.Response
	PostResponseOutput *http.Response
}

func (m MockClient) Get(url string) (resp *http.Response, err error) {
	return m.GetResponseOutput, nil
}

func (m MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return m.PostResponseOutput, nil
}

func TestDoGetRequest(t *testing.T) {
	words := WordsPage{
		Page: Page{"words"},
		Words: Words{
			Input: "This is a test",
			Words: []string{"This", "is", "a", "test"},
		},
	}

	wordsBytes, err := json.Marshal(words)
	if err != nil {
		t.Errorf("Marshal error: %s", err)
	}

	apiInstance := api{
		Options: Options{},
		Client: MockClient{
			GetResponseOutput: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(wordsBytes)),
			},
		},
	}
	response, err := apiInstance.DoGetRequest("http://localhost/words")
	if err != nil {
		t.Errorf("DoGetRequest error: %s", err)
	}
	if response == nil {
		t.Fatalf("Response is empty")
	}
	if response.GetResponse() != strings.Join([]string{"This", "is", "a", "test"}, ", ") {
		t.Errorf("Unexpected response: %s", response.GetResponse())
	}
}
