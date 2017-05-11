package main

import (
	"fury/scraping"
	"fmt"
)


func GlobalHandler(proxy scraping.ScrapingResultProxy) {
	fmt.Println("D")
}

func ScraperHandler(proxy scraping.ScrapingResultProxy) {
	fmt.Println("DWWW")
}


func main() {
	runner := scraping.NewRunner().SetHandler(GlobalHandler)
	scraper := scraping.NewScraper("http://golangweekly.com").SetHandler(ScraperHandler)

	runner.PushScraper(scraper)
	runner.Run()
}