package main


import "fury/scraping"


func main() {
	runner := scraping.NewRunner()
	runner.PushScraper(scraping.NewScraper("golangweekly.com", "http://golangweekly.com"))
	runner.Run()
}