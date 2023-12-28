package net

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type SimpleHttpClient struct {
	httpClient http.Client
}

func NewSimpleHttpClient(httpClient http.Client) *SimpleHttpClient {
	return &SimpleHttpClient{httpClient: httpClient}
}

func (simpleHttpClient SimpleHttpClient) Do(req *http.Request) (*http.Response, error) {
	return simpleHttpClient.httpClient.Do(req)
}
