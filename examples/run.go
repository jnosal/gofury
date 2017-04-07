package main

import (
	"fury"
	"net/url"
)

type TestResource struct {
    fury.DeleteNotSupported
    fury.HeadNotSupported
    fury.PatchNotSupported
    fury.PostNotSupported
    fury.PutNotSupported
fury.TraceNotSupported

}

func (TestResource) Get(values url.Values) (int, interface{}) {
    data := map[string]string{"hello": "world"}
    return 200, data
}


func main() {
	f := fury.New()
	f.AddResource(new(TestResource), "/test")
	f.Start(3000)
}
