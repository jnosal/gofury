package main

import (
	"fmt"
	"fury/scraping"
)

func GlobalHandler(proxy scraping.ScrapingResultProxy) {
	fmt.Println("D")
}

func ScraperHandler(proxy scraping.ScrapingResultProxy) {
	fmt.Println("DWWW")
}

func main() {
	engine := scraping.NewEngine().SetHandler(GlobalHandler)

	engine.PushScraper(scraping.NewScraper("http://golangweekly!!!.com").SetHandler(ScraperHandler))
	engine.PushScraper(scraping.NewScraper("http://golangweekly.com").SetHandler(ScraperHandler))
	engine.Run()
}
