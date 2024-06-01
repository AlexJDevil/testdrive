package testdrive

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arun0009/testdrive/env"
	"github.com/arun0009/testdrive/testdrive/api"
	resty "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type TestServerA struct {
	t assert.TestingT
}

func (t TestServerA) Name() string {
	return "Verify ServerA"
}

func (t TestServerA) Test(client resty.Client, apiResponse chan<- api.Response) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
	}))
	resp, err := client.R().Get(svr.URL)
	var responses []api.ResponseWithError
	responseWithError := api.ResponseWithError{Response: *resp, Error: err}
	responses = append(responses, responseWithError)
	apiResponse <- api.Response{Name: t.Name(), Responses: responses}

}

type TestServerB struct {
	t assert.TestingT
}

func (t TestServerB) Name() string {
	return "Verify ServerB"
}

func (t TestServerB) Test(client resty.Client, apiResponse chan<- api.Response) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	resp, err := client.R().Get(svr.URL)
	var responses []api.ResponseWithError
	responseWithError := api.ResponseWithError{Response: *resp, Error: err}
	responses = append(responses, responseWithError)
	apiResponse <- api.Response{Name: t.Name(), Responses: responses}

}

func TestRun(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	env.ApiURL = svr.URL
	startTime := time.Now()
	td, err := New()
	assert.NoError(t, err)
	td.RegisterApi(&TestServerA{})
	td.RegisterApi(&TestServerB{})
	apiResponses := td.Run()
	assert.Equal(t, 2, len(apiResponses))
	for _, r := range apiResponses {
		for _, responses := range r.Responses {
			assert.NoError(t, responses.Error)
			assert.Regexp(t, "Verify Server.?", r.Name)
			assert.Equal(t, http.StatusOK, responses.Response.StatusCode())
		}
	}
	elapsed := time.Since(startTime)
	assert.Less(t, elapsed.Seconds(), 5.2)
}

func TestRetries(t *testing.T) {
	retryCount := 0
	restyRetryCount := 0
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			return
		}
		if r.URL.Path == "/assets/fail" {
			retryCount++
			w.WriteHeader(http.StatusBadGateway)
			_, _ = fmt.Fprintf(w, "Bad Gateway")
			return
		}
		if r.URL.Path == "/assets/intermittent" {
			retryCount++
			if retryCount < restyRetryCount {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(w, "Internal Server Error")
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, "Done")
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "Success")
	}))
	env.ApiURL = svr.URL
	td, err := New()
	assert.NoError(t, err)

	restyRetryCount = td.Client.RetryCount

	retryCount = 0
	resp, err := td.Client.R().Get(svr.URL + "/health")
	assert.NoError(t, err)
	assert.True(t, resp.IsSuccess())

	retryCount = 0
	resp, err = td.Client.R().Get(svr.URL + "/assets/intermittent")
	assert.NoError(t, err)
	assert.True(t, resp.IsSuccess())
	assert.Equal(t, retryCount, td.Client.RetryCount)

	retryCount = 0
	resp, err = td.Client.R().Get(svr.URL + "/assets/fail")
	assert.NoError(t, err)
	assert.True(t, resp.IsError())
	assert.Greater(t, retryCount, td.Client.RetryCount)
}
