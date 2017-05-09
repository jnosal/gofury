package main

import (
	"fury/scraping"
	"log"
)


func CrawledHandler(proxy scraping.ScrapingResultProxy) {
	log.Println("EXECUTING HANDLER HUEHUEHUEHUEHUEH")
}

func main() {
	runner := scraping.NewRunner(CrawledHandler)
	scraper := scraping.NewScraper("http://golangweekly.com")

	runner.PushScraper(scraper)
	runner.Run()
}