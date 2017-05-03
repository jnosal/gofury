package fury

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
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

const (
	HEADER_CONTENT_TYPE = "Content-Type"
)

const (
	CONTENT_TYPE_JSON = "application/json"
	CONTENT_TYPE_XML  = "application/xml"
)

type Renderer interface {
	Render(code int, name string, data interface{})
	Html(code int, html string)
	String(code int, s string)
	Json(code int, data interface{})
	Xml(code int, data interface{})
	SetContentType(name string)
	WriteWithCode(code int, data []byte)
}

type Meta struct {
	writer  http.ResponseWriter
	request *http.Request
	path    string
	query   url.Values
	headers http.Header
	fury    *Fury
}

func (meta *Meta) SetContentType(name string) {
	meta.writer.Header().Add(HEADER_CONTENT_TYPE, name)
}

func (meta *Meta) WriteWithCode(code int, data []byte) {
	meta.writer.WriteHeader(code)
	meta.writer.Write(data)
}

func (meta *Meta) Render(code int, name string, data interface{}) {
	return
}

func (meta *Meta) Html(code int, html string) {
	return
}

func (meta *Meta) String(code int, s string) {
	meta.WriteWithCode(code, []byte(s))
}

func (meta *Meta) Json(code int, data interface{}) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(meta.writer, err.Error(), http.StatusInternalServerError)
		return
	}

	meta.SetContentType(CONTENT_TYPE_JSON)
	meta.WriteWithCode(code, content)
}

func (meta *Meta) Xml(code int, data interface{}) {
	content, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(meta.writer, err.Error(), http.StatusInternalServerError)
		return
	}

	meta.SetContentType(CONTENT_TYPE_XML)
	meta.WriteWithCode(code, content)
}

type ConnectMixin interface {
	Connect(meta *Meta)
}

type DeleteMixin interface {
	Delete(meta *Meta)
}

type GetMixin interface {
	Get(meta *Meta)
}

type HeadMixin interface {
	Head(meta *Meta)
}

type OptionsMixin interface {
	Options(meta *Meta)
}

type PatchMixin interface {
	Patch(meta *Meta)
}

type PostMixin interface {
	Post(meta *Meta)
}

type PutMixin interface {
	Put(meta *Meta)
}

type TraceMixin interface {
	Trace(meta *Meta)
}

type Fury struct {
	port           int
	host           string
	preMiddleware  []string
	postMiddleware []string
}

func (fury *Fury) Route(path string, resource interface{}) *Fury {
	http.HandleFunc(path, fury.requestHandler(resource))
	return fury
}

func (fury *Fury) UsePre(middleware ...string) *Fury {
	fury.preMiddleware = append(fury.preMiddleware, middleware...)
	return fury
}

func (fury *Fury) UsePost(middleware ...string) *Fury {
	fury.postMiddleware = append(fury.postMiddleware, middleware...)
	return fury
}

func (fury *Fury) Start() {
	address := fmt.Sprintf("%s:%d", fury.host, fury.port)
	log.Printf("STARTING FURY at %s", address)
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
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

		var handler func(m *Meta)

		switch request.Method {
		case CONNECT:
			if resource, ok := resource.(ConnectMixin); ok {
				handler = resource.Connect
			}
		case DELETE:
			if resource, ok := resource.(DeleteMixin); ok {
				handler = resource.Delete
			}
		case GET:
			if resource, ok := resource.(GetMixin); ok {
				handler = resource.Get
			}
		case HEAD:
			if resource, ok := resource.(HeadMixin); ok {
				handler = resource.Head
			}
		case OPTIONS:
			if resource, ok := resource.(OptionsMixin); ok {
				handler = resource.Options
			}
		case POST:
			if resource, ok := resource.(PostMixin); ok {
				handler = resource.Post
			}
		case PUT:
			if resource, ok := resource.(PutMixin); ok {
				handler = resource.Put
			}
		case PATCH:
			if resource, ok := resource.(PatchMixin); ok {
				handler = resource.Patch
			}
		case TRACE:
			if resource, ok := resource.(TraceMixin); ok {
				handler = resource.Trace
			}
		}

		if handler == nil {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var meta = &Meta{writer: rw, request: request, query: request.Form, headers: request.Header, fury: fury}
		handler(meta)
	}
}

func New(host string, port int) *Fury {
	return &Fury{host: host, port: port}
}
