package main

import (
	"fury/scraping"
)


func CrawledHandler(proxy scraping.ScrapingResultProxy) {

}

func main() {
	runner := scraping.NewRunner(CrawledHandler)
	scraper := scraping.NewScraper("http://golangweeklyddd.com")

	runner.PushScraper(scraper)
	runner.Run()
}