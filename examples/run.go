package main

import (
	"fury"
	"net/url"
	"net/http"
)

type TestResource struct {
	fury.GetSupported

}

func (TestResource) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {
    data := map[string]string{"hello": "world"}
    return 200, data, http.Header{"Content-type": {"application/json"}}
}


func main() {
	f := fury.New()
	f.AddResource(new(TestResource), "/test")
	f.Start(3000)
}
