package main

import (
	"fury"
)

type TestResource struct {
	fury.GetMixin
}

func (TestResource) Get(meta *fury.Meta) {
	data := map[string]string{"hello": "world"}
	meta.Json(200, data)
}

func main() {
	f := fury.New("localhost", 3000)
	f.Route("/test", new(TestResource)).
		Route("/test2", new(TestResource)).
		Route("/test3", new(TestResource))
	f.Start()
}
