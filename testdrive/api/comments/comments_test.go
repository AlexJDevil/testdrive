package comments

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arun0009/testdrive/env"
	resty "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestCommentsAPIHealth(t *testing.T) {
	var tests = []struct {
		statusCode    int
		statusMessage string
		err           bool
		setURL        bool
	}{
		{statusCode: http.StatusOK, statusMessage: http.StatusText(http.StatusOK), err: false, setURL: true},
		{statusCode: http.StatusInternalServerError, statusMessage: http.StatusText(http.StatusInternalServerError), err: false, setURL: true},
		{statusCode: 0, statusMessage: "", err: true, setURL: false},
	}
	for _, tt := range tests {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.statusCode)
			_, _ = fmt.Fprintf(w, tt.statusMessage)
		}))
		if tt.setURL == true {
			env.ApiURL = svr.URL
		}
		client := resty.New()
		response, err := health(*client)
		assert.Equal(t, tt.statusCode, response.StatusCode())
		if tt.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		env.ApiURL = ""
	}
}

func TestCommentsAPIGetEndpoint(t *testing.T) {
	var tests = []struct {
		statusCode    int
		statusMessage string
		err           bool
		setURL        bool
	}{
		{statusCode: http.StatusOK, statusMessage: http.StatusText(http.StatusOK), err: false, setURL: true},
		{statusCode: http.StatusInternalServerError, statusMessage: http.StatusText(http.StatusInternalServerError), err: false, setURL: true},
		{statusCode: 0, statusMessage: "", err: true, setURL: false},
	}
	for _, tt := range tests {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.statusCode)
			w.Write([]byte(`
					[
						{
							"id": 1,
							"body": "some comment",
							"postId": 1
						},
						{
							"id": 2,
							"body": "some comment",
							"postId": 1
						}
					]
			`))
			_, _ = fmt.Fprintf(w, tt.statusMessage)
		}))
		if tt.setURL == true {
			env.ApiURL = svr.URL
		}
		client := resty.New()
		response, err := getComments(*client)
		assert.Equal(t, tt.statusCode, response.StatusCode())
		if tt.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		env.ApiURL = ""
	}
}

func TestCommentsAPIPostEndpoint(t *testing.T) {
	var tests = []struct {
		statusCode    int
		statusMessage string
		err           bool
		setURL        bool
	}{
		{statusCode: http.StatusOK, statusMessage: `{
			"posts": [
			  { "id": 1, "title": "json-server", "author": "typicode" }
			]
		  }`, err: false, setURL: true},
		{statusCode: http.StatusInternalServerError, statusMessage: http.StatusText(http.StatusInternalServerError), err: false, setURL: true},
		{statusCode: 0, statusMessage: "", err: true, setURL: false},
	}
	for _, tt := range tests {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.statusCode)
			_, _ = fmt.Fprintf(w, tt.statusMessage)
		}))
		if tt.setURL == true {
			env.ApiURL = svr.URL
		}
		client := resty.New()
		response, err := postComments(*client)
		assert.Equal(t, tt.statusCode, response.StatusCode())
		if tt.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		env.ApiURL = ""
	}
}
