package scraping

import (
	"fury"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	URL "net/url"
	"strings"
	"sync"
	"time"
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

type ScrapingResultProxy struct {
	resp    http.Response
	scraper Scraper
}

func (proxy ScrapingResultProxy) String() (result string) {
	result = "Twoja stara"
	return
}

type Runnable interface {
	Run() (err error)
}

type DefaultRunner struct {
	limitCrawl int
	limitFail int
	handler func(ScrapingResultProxy)
	finished  int
	scrapers  []*Scraper
	chDone    chan *Scraper
	chScraped chan ScrapingResultProxy
}

func (runner *DefaultRunner) IncrFinishedCounter() {
	runner.finished += 1
}

func (runner DefaultRunner) Done() bool {
	return len(runner.scrapers) == runner.finished
}

func (runner *DefaultRunner) Run() {
	defer runner.Close()

	for _, scraper := range runner.scrapers {
		go scraper.Start(scraper.baseUrl)
	}

	// main scraping loop
	for {
		select {
		case proxy, ok := <-runner.chScraped:
			if !ok {
				break
			}
			runner.handler(proxy)
		case scraper, ok := <-runner.chDone:
			fury.Logger().Infof("Stopped %s", scraper)
			runner.IncrFinishedCounter()
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
	close(runner.chScraped)
}

func (runner *DefaultRunner) PushScraper(scrapers ...*Scraper) *DefaultRunner {
	for _, scraper := range scrapers {
		fury.Logger().Debugf("Attaching new scraper %s", scraper)
		scraper.runner = runner
	}
	runner.scrapers = append(runner.scrapers, scrapers...)
	return runner
}

type Scraper struct {
	crawled      int
	successful   int
	failed       int
	fetchMutex   *sync.Mutex
	crawledMutex *sync.Mutex
	domain       string
	baseUrl      string
	fetchedUrls  map[string]bool
	runner       *DefaultRunner
	extractor    *LinkExtractor
}

func (scraper *Scraper) IncrCounters(isSuccessful bool) {
	scraper.crawledMutex.Lock()
	scraper.crawled += 1
	if isSuccessful {
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
	if scraper.crawled == scraper.runner.limitCrawl {
		fury.Logger().Warningf("Crawl limit exceeded: %s", scraper)
		ok = true
	} else if scraper.failed == scraper.runner.limitFail {
		fury.Logger().Warningf("Fail limit exceeeded: %s", scraper)
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
			scraper.Fetch(url, true)
		}

		if scraper.CheckIfShouldStop() {
			scraper.Stop()
		}
	})
}

func (scraper *Scraper) Stop() {
	fury.Logger().Infof("Stopping %s", scraper)
	scraper.runner.chDone <- scraper
}

func (scraper *Scraper) Start(baseUrl string) {
	fury.Logger().Infof("Starting: %s", scraper)

	resp, err := scraper.Fetch(baseUrl, false)

	if err != nil {
		fury.Logger().Errorf("Base url is corrupted %s", baseUrl)
		scraper.Stop()
		return
	}

	scraper.RunExtractor(resp)

	return
}

func (scraper *Scraper) Notify(resp *http.Response) {
	scraper.runner.chScraped <- NewResultProxy(*scraper, *resp)
}

func (scraper *Scraper) Fetch(url string, extract bool) (resp *http.Response, err error) {
	if ok := scraper.CheckIfFetched(url); ok {
		return
	}
	scraper.MarkAsFetched(url)

	fury.Logger().Infof("Fetching: %s", url)
	tic := time.Now()

	resp, err = http.Get(url)

	fury.Logger().Debugf("Request to %s took: %s", url, time.Since(tic))

	scraper.IncrCounters(err == nil)
	if err != nil {
		fury.Logger().Warningf("Failed to crawl %s", url)
		return
	}

	if extract {
		scraper.RunExtractor(resp)
	}

	scraper.Notify(resp)
	return
}

func (scraper *Scraper) String() (result string) {
	result = fmt.Sprintf("<Scraper: %s>. Crawled: %d, successful: %d failed: %d.",
		scraper.domain, scraper.crawled, scraper.successful, scraper.failed)
	return
}

func NewRunner(handler func(ScrapingResultProxy)) (r *DefaultRunner) {
	r = &DefaultRunner{
		limitCrawl: 1000,
		limitFail: 50,
		handler: handler,
		finished:  0,
		chDone:    make(chan *Scraper),
		chScraped: make(chan ScrapingResultProxy),
	}
	return
}

func NewScraper(sourceUrl string) (s *Scraper) {
	parsed, err := URL.Parse(sourceUrl)
	if err != nil {
		fury.Logger().Infof("Inappropriate URL: %s", sourceUrl)
		return
	}
	s = &Scraper{
		crawled:      0,
		successful:   0,
		failed:       0,
		domain:       parsed.Host,
		baseUrl:      sourceUrl,
		fetchedUrls:  make(map[string]bool),
		crawledMutex: &sync.Mutex{},
		fetchMutex:   &sync.Mutex{},
		extractor:    &LinkExtractor{},
	}
	return
}

func NewResultProxy(scraper Scraper, resp http.Response) (result ScrapingResultProxy) {
	result = ScrapingResultProxy{scraper: scraper, resp: resp}
	return
}
