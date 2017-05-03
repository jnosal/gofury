package main

import (
	"fury"
)

type JsonResource struct {
	fury.GetMixin
}

func (JsonResource) Get(meta *fury.Meta) {
	data := map[string]string{"hello": "world"}
	meta.Json(200, data)
}

type StringResource struct {
	fury.GetMixin
}

func (StringResource) Get(meta *fury.Meta) {
	meta.String(200, "test")
}

func main() {
	f := fury.New("localhost", 3000)
	f.Route("/test", new(JsonResource)).
		Route("/test2", new(StringResource))
	f.Start()
}
