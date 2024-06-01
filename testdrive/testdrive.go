package testdrive

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/arun0009/testdrive/alert"
	"github.com/arun0009/testdrive/testdrive/api"
	"github.com/arun0009/testdrive/testdrive/api/comments"
	resty "github.com/go-resty/resty/v2"
)

type (
	TestDrive struct {
		Client   resty.Client
		Notifier alert.Notifier
		TestApi  []api.Api
	}
)

func New() (*TestDrive, error) {
	client := resty.New().SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: false,
	})
	//set auth token etc here as global headers for all requests
	//client = client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))

	// Retries are configured per client
	client = client.SetRetryCount(1).SetRetryWaitTime(3 * time.Second).SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.IsError()
		})

	slackNoifier := new(alert.SlackNotifer)

	return &TestDrive{*client, slackNoifier, nil}, nil
}

func (t *TestDrive) RegisterApi(api api.Api) {
	t.TestApi = append(t.TestApi, api)
}

func (t *TestDrive) Run() []api.Response {
	ch := make(chan api.Response)
	var apiResponses []api.Response
	for _, api := range t.TestApi {
		go api.Test(t.Client, ch)
	}
	for range t.TestApi {
		apiResponses = append(apiResponses, <-ch)
	}

	for _, apiResponse := range apiResponses {
		var messages []string
		for _, resp := range apiResponse.Responses {
			if resp.Error != nil {
				messages = append(messages, fmt.Sprintf("%s: %s", apiResponse.Name, resp.Error.Error()))
			}
			if !resp.Response.IsSuccess() {
				messages = append(messages, fmt.Sprintf("%s: %d", apiResponse.Name, resp.Response.StatusCode()))
			}
		}
		if len(messages) > 0 {
			log.Printf("error message is %s", strings.Join(messages, "\n"))
			t.Notifier.Notify(strings.Join(messages, "\n"))
		}
	}
	return apiResponses
}

func Exec() {
	td, err := New()
	if err != nil {
		td.Notifier.Notify(err.Error())
	}
	td.RegisterApi(comments.Api{})
	td.Run()
}
