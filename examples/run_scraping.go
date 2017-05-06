package main


import "log"
import "fury/scraping"


func main() {
	runner := scraping.NewRunner()
	runner.PushScraper(scraping.NewScraper("golangweekly.com", "http://golangweekly.com"))

	runner.Run()
	defer runner.Close()


	for {
		select {
		case url, ok := <-runner.ChFoundUrls:
			if !ok {
				break
			}
			log.Print(url)
		case _, ok := <- runner.ChDone:
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