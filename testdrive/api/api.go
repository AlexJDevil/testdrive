package api

import (
	resty "github.com/go-resty/resty/v2"
)

type Api interface {
	Name() string
	Test(client resty.Client, response chan<- Response)
}

type ResponseWithError struct {
	Response resty.Response
	Error    error
}

type Response struct {
	Name      string
	Responses []ResponseWithError
}
