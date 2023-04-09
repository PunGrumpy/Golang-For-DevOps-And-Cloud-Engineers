package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type MockRoundTripper struct {
	RoundTripperOutput *http.Response
}

func (m MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Authorization") != "Bearer abc" {
		return nil, fmt.Errorf("Wrong authorization header: %s", req.Header.Get("Authorization"))
	}
	return m.RoundTripperOutput, nil
}

func TestRoundTrip(t *testing.T) {
	loginResponse := LoginResponse{
		Token: "abc",
	}
	loginResponseBytes, err := json.Marshal(loginResponse)
	if err != nil {
		t.Errorf("Marshal error: %s", err)
	}
	myJWTTransport := MyJWTTransport{
		transport: MockRoundTripper{
			RoundTripperOutput: &http.Response{
				StatusCode: 200,
			},
		},
		HTTPClient: MockClient{
			PostResponseOutput: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(loginResponseBytes)),
			},
		},
		password: "xyz",
	}
	req := &http.Request{
		Header: make(http.Header),
	}
	res, err := myJWTTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Status code is not 200: %d", res.StatusCode)
	}
}
