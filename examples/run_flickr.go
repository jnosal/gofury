package main

import (
	"flag"
	"fmt"
	"fury"
	"io/ioutil"
	"net/http"
)

type FlickrConfig struct {
	ApiKey    string
	ApiSecret string
	BaseUrl   string
}

type FlickrClient struct {
	config *FlickrConfig
}

func (client *FlickrClient) Run() {
	res, err := http.Get(client.config.BaseUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(body)
}

func (client FlickrClient) String() (s string) {
	s = "Flickr API Client"
	return s
}

func NewConfig(key string, secret string) *FlickrConfig {
	return &FlickrConfig{
		ApiKey:    key,
		ApiSecret: secret,
		BaseUrl:   "https://api.flickr.com/services/rest/",
	}
}

func (config FlickrConfig) String() (s string) {
	s = fmt.Sprintf("Flickr Config. Key: %s Secret: %s", config.ApiKey, config.ApiSecret)
	return
}

type SearchResource struct {
}

func (resource *SearchResource) Get(meta *fury.Meta) {
	client := meta.App().FromRegistry("flickr").(FlickrClient)
	client.Run()
	//data := map[string]string{"hello": "world"}
	meta.Json(http.StatusOK, nil)
}

func main() {
	key := flag.String("api-key", "", "Flickr API KEY")
	secret := flag.String("api-secret", "", "Flickr API SECRET")
	flag.Parse()

	config := NewConfig(*key, *secret)
	client := FlickrClient{config}

	fury.Logger().Infof(config.String())

	f := fury.New("localhost", 3000).
		UseMiddleware(fury.RequestStatsMiddleware).
		ToRegistry("flickr", client).
		Route("/search", new(SearchResource))

	f.Start()
}
