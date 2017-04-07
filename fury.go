package fury

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	DELETE = "DELETE"
	GET    = "GET"
	HEAD   = "HEAD"
	PATCH  = "PATCH"
	POST   = "POST"
	PUT    = "PUT"
	TRACE  = "TRACE"
)

type Resource interface {
	Delete(values url.Values) (int, interface{})
	Get(values url.Values) (int, interface{})
	Head(values url.Values) (int, interface{})
	Patch(values url.Values) (int, interface{})
	Post(values url.Values) (int, interface{})
	Put(values url.Values) (int, interface{})
	Trace(values url.Values) (int, interface{})
}

type (
	DeleteNotSupported struct{}
	GetNotSupported    struct{}
	HeadNotSupported   struct{}
	PatchNotSupported  struct{}
	PostNotSupported   struct{}
	PutNotSupported    struct{}
	TraceNotSupported  struct{}
)

func (DeleteNotSupported) Delete(values url.Values) (int, interface{}) {
	return 405, ""
}

func (GetNotSupported) Get(values url.Values) (int, interface{}) {
	return 405, ""
}

func (HeadNotSupported) Head(values url.Values) (int, interface{}) {
	return 405, ""
}

func (PatchNotSupported) Patch(values url.Values) (int, interface{}) {
	return 405, ""
}

func (PostNotSupported) Post(values url.Values) (int, interface{}) {
	return 405, ""
}

func (PutNotSupported) Put(values url.Values) (int, interface{}) {
	return 405, ""
}

func (TraceNotSupported) Trace(values url.Values) (int, interface{}) {
	return 405, ""
}

type Fury struct {
	Name string
}

func (fury *Fury) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

func (fury *Fury) requestHandler(resource Resource) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var data interface{}
		var code int

		request.ParseForm()
		method := request.Method
		values := request.Form

		switch method {
		case DELETE:
			code, data = resource.Delete(values)
		case GET:
			code, data = resource.Get(values)
		case HEAD:
			code, data = resource.Head(values)
		case PATCH:
			code, data = resource.Patch(values)
		case POST:
			code, data = resource.Post(values)
		case PUT:
			code, data = resource.Put(values)
		case TRACE:
			code, data = resource.Trace(values)
		default:
			fury.Abort(rw, 405)
			return
		}

		content, err := json.Marshal(data)
		if err != nil {
			fury.Abort(rw, 500)
		}
		rw.WriteHeader(code)
		rw.Write(content)
	}
}

func (fury *Fury) AddResource(resource Resource, path string) {
	http.HandleFunc(path, fury.requestHandler(resource))
}

func (fury *Fury) Start(port int) {
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, nil)
}

func New() (f *Fury) {
	f = &Fury{
		Name: "Hello Fury!!!",
	}
	return
}
