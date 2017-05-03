package main

import (
	"fury"
	"net/http"
)

type JsonResource struct {
	fury.GetMixin
}

func (JsonResource) Get(meta *fury.Meta) {
	data := map[string]string{"hello": "world"}
	meta.Json(http.StatusOK, data)
}

type StringResource struct {
	fury.GetMixin
}

func (StringResource) Get(meta *fury.Meta) {
	meta.String(http.StatusOK, "test")
}

func main() {
	f := fury.New("localhost", 3000)
	f.Route("/test", new(JsonResource)).
		Route("/test2", new(StringResource))
	f.Start()
}
