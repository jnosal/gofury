package main

import (
	"fury"
	"log"
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

type DetailResource struct {
}

func (resource *DetailResource) Get(meta *fury.Meta) {
	log.Println(meta.RequestHeaders())
	fury.RetrieveResource(resource, meta)
}

func (resource *DetailResource) Retrieve() (s string, err error) {
	s = "DDDCC!!!!!!!!1321312333333444"
	return
}

func main() {
	f := fury.New("localhost", 3000)
	f.Route("/test", new(JsonResource)).
		Route("/test2", new(StringResource)).
		Route("/detail", new(DetailResource))
	f.Start()
}
