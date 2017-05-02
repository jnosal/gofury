package fury

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

type ConnectSupported interface {
	Connect(url.Values, http.Header) (int, interface{}, http.Header)
}

type DeleteSupported interface {
	Delete(url.Values, http.Header) (int, interface{}, http.Header)
}

type GetSupported interface {
	Get(url.Values, http.Header) (int, interface{}, http.Header)
}

type HeadSupported interface {
	Head(url.Values, http.Header) (int, interface{}, http.Header)
}

type OptionsSupported interface {
	Options(url.Values, http.Header) (int, interface{}, http.Header)
}

type PatchSupported interface {
	Patch(url.Values, http.Header) (int, interface{}, http.Header)
}

type PostSupported interface {
	Post(url.Values, http.Header) (int, interface{}, http.Header)
}

type PutSupported interface {
	Put(url.Values, http.Header) (int, interface{}, http.Header)
}

type TraceSupported interface {
	Trace(url.Values, http.Header) (int, interface{}, http.Header)
}

type Fury struct {
	port       int
	host       string
	middleware []string
}

func (fury *Fury) Route(path string, resource interface{}) *Fury {
	http.HandleFunc(path, fury.requestHandler(resource))
	return fury
}

func (fury *Fury) Use(middleware ...string) *Fury {
	fury.middleware = append(fury.middleware, middleware...)
	return fury
}

func (fury *Fury) Start() error {
	address := fmt.Sprintf("%s:%d", fury.host, fury.port)
	log.Printf("STARTING FURY at %s", address)
	return http.ListenAndServe(address, nil)
}

func (fury *Fury) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

func (fury *Fury) requestHandler(resource interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		if request.ParseForm() != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var handler func(url.Values, http.Header) (int, interface{}, http.Header)

		switch request.Method {
		case CONNECT:
			if resource, ok := resource.(ConnectSupported); ok {
				handler = resource.Connect
			}
		case DELETE:
			if resource, ok := resource.(DeleteSupported); ok {
				handler = resource.Delete
			}
		case GET:
			if resource, ok := resource.(GetSupported); ok {
				handler = resource.Get
			}
		case HEAD:
			if resource, ok := resource.(HeadSupported); ok {
				handler = resource.Head
			}
		case OPTIONS:
			if resource, ok := resource.(OptionsSupported); ok {
				handler = resource.Options
			}
		case POST:
			if resource, ok := resource.(PostSupported); ok {
				handler = resource.Post
			}
		case PUT:
			if resource, ok := resource.(PutSupported); ok {
				handler = resource.Put
			}
		case PATCH:
			if resource, ok := resource.(PatchSupported); ok {
				handler = resource.Patch
			}
		case TRACE:
			if resource, ok := resource.(TraceSupported); ok {
				handler = resource.Trace
			}
		}

		if handler == nil {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		code, data, header := handler(request.Form, request.Header)

		content, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		for name, values := range header {
			for _, value := range values {
				rw.Header().Add(name, value)
			}
		}
		rw.WriteHeader(code)
		rw.Write(content)
	}
}

func New(host string, port int) *Fury {
	return &Fury{host: host, port: port}
}
