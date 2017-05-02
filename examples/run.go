package main

import (
	"fury"
	"net/http"
	"net/url"
)

type TestResource struct {
	fury.GetSupported
}

func (TestResource) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {
	data := map[string]string{"hello": "world"}
	return 200, data, http.Header{"Content-type": {"application/json"}}
}

func main() {
	f := fury.New("localhost", 3000)
	f.Route("/test", new(TestResource)).
		Route("/test2", new(TestResource)).
		Route("/test3", new(TestResource))
	f.Start()
}
