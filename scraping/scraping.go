package scraping

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const HREF = "href"
const LIMIT_CRAWL = 10
const LIMIT_FAIL = 1

func GetHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == HREF {
			href = a.Val
			ok = true
		}
	}

	return
}

type Extractable interface {
	Extract()
}

type LinkExtractor struct {
	Extractable
}

func (extractor *LinkExtractor) Extract(r io.ReadCloser, callback func(string)) {
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
			callback(url)
		}
	}
}

type Runnable interface {
	Run() (err error)
}

type DefaultRunner struct {
	finished    int
	scrapers    []*Scraper
	chDone      chan *Scraper
	chFoundUrls chan string
}

func (runner *DefaultRunner) FinishedIncr() {
	runner.finished += 1
}

func (runner DefaultRunner) Done() bool {
	return len(runner.scrapers) == runner.finished
}

func (runner *DefaultRunner) Run() {
	defer runner.Close()

	for _, scraper := range runner.scrapers {
		log.Printf("Starting: %s", scraper)
		go scraper.Start(scraper.baseUrl)
	}

	// main scraping loop
	for {
		select {
		case _, ok := <-runner.chFoundUrls:
			if !ok {
				break
			}
		case scraper, ok := <-runner.chDone:
			log.Printf("%s is done.", scraper)
			runner.FinishedIncr()
			if !ok {
				break
			}

		}
		if runner.Done() {
			break
		}
	}
}

func (runner *DefaultRunner) Close() {
	close(runner.chDone)
	close(runner.chFoundUrls)
}

func (runner *DefaultRunner) PushScraper(scrapers ...*Scraper) *DefaultRunner {
	for _, scraper := range scrapers {
		scraper.runner = runner
	}
	runner.scrapers = append(runner.scrapers, scrapers...)
	return runner
}

type Scraper struct {
	crawled int
	successful int
	failed int
	fetchMutex  *sync.Mutex
	crawledMutex *sync.Mutex
	domain      string
	baseUrl     string
	fetchedUrls map[string]bool
	runner      *DefaultRunner
	extractor   *LinkExtractor
}


func (scraper *Scraper) IncrCounters(isSuccessful bool) {
	scraper.crawledMutex.Lock()
	scraper.crawled += 1
	if (isSuccessful) {
		scraper.successful += 1
	} else {
		scraper.failed += 1
	}
	scraper.crawledMutex.Unlock()
}


func (scraper *Scraper) MarkAsFetched(url string) {
	scraper.fetchMutex.Lock()
	scraper.fetchedUrls[url] = true
	scraper.fetchMutex.Unlock()
}


func (scraper *Scraper) CheckIfShouldStop() (ok bool) {
	scraper.crawledMutex.Lock()
	if scraper.crawled == LIMIT_CRAWL {
		log.Printf("Crawl limit exceeded: %s", scraper)
		ok = true
	} else if scraper.failed == LIMIT_FAIL {
		log.Printf("Fail limit exceeeded: %s", scraper)
		ok = true
	}
	scraper.crawledMutex.Unlock()
	return
}

func (scraper *Scraper) CheckIfFetched(url string) (ok bool) {
	scraper.fetchMutex.Lock()
	_, ok = scraper.fetchedUrls[url]
	scraper.fetchMutex.Unlock()
	return
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


func (scraper *Scraper) RunExtractor(resp *http.Response) {
	scraper.extractor.Extract(resp.Body, func(url string) {
		ok, url := scraper.CheckUrl(url)

		if ok {
			go scraper.Fetch(url, true)
		}

		if scraper.CheckIfShouldStop() {
			scraper.Stop()
		}
	})
}


func (scraper *Scraper) Stop() {
	scraper.runner.chDone <- scraper
}


func (scraper *Scraper) Start(baseUrl string) {
	resp, err :=  scraper.Fetch(baseUrl, false)

	if err != nil {
		log.Printf("Base url is corrupted stopping %s", scraper)
		scraper.Stop()
		return
	}

	scraper.RunExtractor(resp)

	return
}


func (scraper *Scraper) Fetch(url string, extract bool) (resp *http.Response, err error) {
	if ok := scraper.CheckIfFetched(url); ok {
		return
	}
	scraper.MarkAsFetched(url)

	log.Printf("Fetching: %s", url)
	resp, err = http.Get(url)

	scraper.IncrCounters(err == nil)
	if err != nil {
		log.Printf("ERROR: Failed to crawl %s", url)
		return
	}

	if extract {
		scraper.RunExtractor(resp)
	}

	return
}

func (scraper *Scraper) String() (result string) {
	result = fmt.Sprintf("<Scraper: %s>. Crawled: %d, successful: %d failed: %d.", scraper.domain, scraper.crawled, scraper.successful, scraper.failed)
	return
}

func NewRunner() (r *DefaultRunner) {
	r = &DefaultRunner{
		finished:    0,
		chDone:      make(chan *Scraper),
		chFoundUrls: make(chan string),
	}
	return
}

func NewScraper(domain string, url string) (s *Scraper) {
	s = &Scraper{
		crawled: 0,
		successful: 0,
		failed: 0,
		domain:      domain,
		baseUrl:     url,
		fetchedUrls: make(map[string]bool),
		crawledMutex: &sync.Mutex{},
		fetchMutex:  &sync.Mutex{},
		extractor:   &LinkExtractor{},
	}
	return
}
