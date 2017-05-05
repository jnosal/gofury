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


type DetailResource struct {
}


func (d *DetailResource) Get(meta *fury.Meta) {
	fury.RetrieveResource(d, meta)
}

func (g *DetailResource) Retrieve() string {
	return "DDDCC!!!!!!!!1321312333333444"
}


func main() {
	f := fury.New("localhost", 3000)
	f.Route("/test", new(JsonResource)).
		Route("/test2", new(StringResource)).
		Route("/detail", new(DetailResource))
	f.Start()
}
