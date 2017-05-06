package scraping

import (
	"log"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"io"
	"sync"
)


func GetHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	return
}


type Runnable interface {
	Run() (err error)
}

type DefaultRunner struct {
	finished int
	scrapers []*Scraper
	ChDone chan bool
	ChFoundUrls chan string
}


func NewRunner() (r *DefaultRunner) {
	r = &DefaultRunner{
		finished: 0,
		ChDone:    make(chan bool),
		ChFoundUrls: make(chan string),
	}
	return
}


func NewScraper(domain string, url string) (s *Scraper) {
	s = &Scraper{
		domain: domain,
		baseUrl: url,
		fetchedUrls: make(map[string]bool),
		mutex: &sync.Mutex{},
	}
	return
}

func (runner *DefaultRunner) FinishedIncr() {
	runner.finished += 1

}

func (runner DefaultRunner) Done() bool {
	return len(runner.scrapers) == runner.finished
}

func (runner *DefaultRunner) Run() {
	for _, scraper := range runner.scrapers {
		log.Printf("Starting scraper: %s", scraper.baseUrl)
		go scraper.Fetch(scraper.baseUrl)
	}
}

func (runner *DefaultRunner) Close() {
	close(runner.ChDone)
	close(runner.ChFoundUrls)
}


func (runner *DefaultRunner) PushScraper(scrapers ...*Scraper) *DefaultRunner {
	for _, scraper := range scrapers {
		scraper.runner = runner
	}
	runner.scrapers = append(runner.scrapers, scrapers...)
	return runner
}


type Scraper struct {
	mutex  *sync.Mutex
	domain string
	baseUrl string
	fetchedUrls map[string]bool
	runner *DefaultRunner
}


func (scraper *Scraper) CheckUrl(sourceUrl string) (ok bool, url string) {
	if strings.Contains(sourceUrl, scraper.domain) && strings.Index(sourceUrl, "http") == 0 {
		url = sourceUrl
		ok = true
	} else if strings.Index(sourceUrl, "/") == 0 {
		url = scraper.baseUrl + sourceUrl
		ok = true
	}
	return
}

func (scraper *Scraper) ExtractLinks(r io.ReadCloser) {
	z := html.NewTokenizer(r)
	defer r.Close()

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}


			ok, url := GetHref(t)
			if !ok {
				continue
			}

			ok, url = scraper.CheckUrl(url)

			if ok {
				go scraper.Fetch(url)
			}
		}
	}
}


func (scraper *Scraper) Fetch(url string) {
	scraper.mutex.Lock()
	_, ok := scraper.fetchedUrls[url]
	scraper.mutex.Unlock()

	if ok {
		return
	}

	scraper.mutex.Lock()
	scraper.fetchedUrls[url] = true
	scraper.mutex.Unlock()

	log.Printf("Fetching: %s", url)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("ERROR: Failed to crawl %s", url)
		return
	}

	go scraper.ExtractLinks(resp.Body)
	return
}