package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/silver-xu/gotask"
)

func main() {
	start := time.Now()
	stockSymbols := []string{
		"AAL",
		"AAPL",
		"ADBE",
		"ADI",
		"ADP",
		"ADSK",
		"ALGN",
		"ALXN",
		"AMAT",
		"AMGN",
		"AMZN",
		"ATVI",
		"ASML",
		"AVGO",
		"BIDU",
		"BIIB",
		"BMRN",
		"CA",
		"CDNS",
		"CELG",
	}

	urlPattern := "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=TN1GWKC8BTNTPZ7P"
	jobs := make(map[string]func() (interface{}, error))

	for _, symbol := range stockSymbols {
		url := fmt.Sprintf(urlPattern, symbol)

		jobs[symbol] = func() (interface{}, error) {
			resp, err := http.Get(url)

			if err == nil {
				defer resp.Body.Close()
			} else {
				return nil, err
			}

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				return nil, err
			}

			return string(body), err
		}
	}

	results, errs := gotask.WhenAll(jobs, 10)

	for key, ret := range results {
		fmt.Println("key " + key + " has result of: " + ret.(string))
	}

	for key, err := range errs {
		fmt.Println("key " + key + " has error of: " + err.(error).Error())
	}

	elapsed := time.Since(start)
	log.Printf("Stock pulling took %s", elapsed)
}
