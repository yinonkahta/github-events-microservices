package main

import "net/http"

type GitHubEventsApi interface {
	List(http.ResponseWriter, *http.Request)
	Count(http.ResponseWriter, *http.Request)
	KRecent(http.ResponseWriter, *http.Request)
}
