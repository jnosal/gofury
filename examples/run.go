package main

import (
	"errors"
	"fmt"
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

type sampleError string

func (e sampleError) Error() string { return "Request: " + string(e) }

type DetailResource struct {
}

func (resource *DetailResource) Delete(meta *fury.Meta) {
	fury.RemoveResource(resource, meta)
}

func (resource *DetailResource) Remove() (s string, err error) {
	s = "NOT_GENERIC"
	return
}

func AnotherMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		fmt.Println("Another")
		next(rw, request)
		fmt.Println("Another 2")
	}
}

func SimpleMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		fmt.Println("Simple")
		next(rw, request)
		fmt.Println("Simple 2")
	}
}

type Sample struct {
	Name string
}

func (s Sample) OK() error {
	return errors.New("Test")
}

func main() {
	f := fury.New("localhost", 3000)
	f.UseMiddleware(fury.RequestStatsMiddleware).
		RegisterCleanup(func(f *fury.Fury) {
			fmt.Println(f)
		}).
		UseMiddleware(fury.RequestCIDMiddleware).
		Route("/test", new(JsonResource)).
		Route("/test2", new(StringResource)).
		Route("/detail", new(DetailResource))
	f.Start()
}
