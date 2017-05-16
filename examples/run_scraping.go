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
	engine := scraping.NewEngine().SetHandler(GlobalHandler)
	scraper := scraping.NewScraper("http://golangweekly.com").SetHandler(ScraperHandler)

	engine.PushScraper(scraper)
	engine.Run()
}