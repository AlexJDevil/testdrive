package comments

import (
	"fmt"

	"github.com/arun0009/testdrive/env"
	"github.com/arun0009/testdrive/testdrive/api"
	resty "github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type Api struct {
}

func (a Api) Name() string {
	return "Comments API"
}

var headers = map[string]string{
	"Accept":       "application/json;v=1",
	"Content-Type": "application/json;v=1",
}

func (a Api) Test(client resty.Client, apiResponse chan<- api.Response) {
	var responses []api.ResponseWithError
	//add each endpoint response to responses
	resp, err := health(client)
	responses = append(responses, api.ResponseWithError{Response: *resp, Error: err})
	resp, err = getComments(client)
	responses = append(responses, api.ResponseWithError{Response: *resp, Error: err})
	resp, err = postComments(client)
	responses = append(responses, api.ResponseWithError{Response: *resp, Error: err})
	apiResponse <- api.Response{Name: a.Name(), Responses: responses}
}

func health(client resty.Client) (*resty.Response, error) {
	resp, err := client.R().SetHeaders(headers).Get(env.ApiURL)
	return resp, err
}

func getComments(client resty.Client) (*resty.Response, error) {
	resp, err := client.R().SetHeaders(headers).Get(env.ApiURL + "/comments")
	if err != nil || resp.IsError() {
		return resp, err
	}
	json := string(resp.Body())
	if gjson.Get(json, "0.id").Num != 1 {
		err = fmt.Errorf("expected first comment id 1 not found in response")
	}
	return resp, err
}

func postComments(client resty.Client) (*resty.Response, error) {
	comment := `{
		"posts": [
		  { "id": 4, "title": "testdrive post", "author": "arun0009" }
		]
	  }`
	resp, err := client.R().SetHeaders(headers).SetBody(comment).Post(env.ApiURL + "/comments")
	if err != nil || resp.IsError() {
		return resp, err
	}
	return resp, err
}
