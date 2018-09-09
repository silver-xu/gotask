package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/silver-xu/gotask"
)

func main() {

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
		"CERN",
		"CHKP",
		"CHTR",
		"CTRP",
		"CTAS",
		"CSCO",
		"CTXS",
		"CMCSA",
		"COST",
		"CSX",
		"CTSH",
		"DISH",
		"DLTR",
		"EA",
		"EBAY",
		"ESRX",
		"EXPE",
		"FAST",
		"FB",
		"FISV",
		"FOX",
		"FOXA",
		"GILD",
		"GOOG",
		"GOOGL",
		"HAS",
		"HSIC",
		"HOLX",
		"ILMN",
		"INCY",
		"INTC",
		"INTU",
		"ISRG",
		"Symbol",
		"IDXX",
		"JBHT",
		"JD",
		"KLAC",
		"KHC",
		"LBTYA",
		"LBTYK",
		"LRCX",
		"MELI",
		"MAR",
		"MCHP",
		"MDLZ",
		"MNST",
		"MSFT",
		"MU",
		"MXIM",
		"MYL",
		"NFLX",
		"NTES",
		"NVDA",
		"ORLY",
		"PAYX",
		"PCAR",
		"BKNG",
		"PYPL",
		"QCOM",
		"QRTEA",
		"REGN",
		"ROST",
		"STX",
		"SHPG",
		"SIRI",
		"SWKS",
		"SBUX",
		"SYMC",
		"SNPS",
		"TTWO",
		"TSLA",
		"TXN",
		"TMUS",
		"ULTA",
		"VOD",
		"VRTX",
		"WBA",
		"WDC",
		"WDAY",
		"XRAY",
		"VRSK",
		"WYNN",
		"XLNX",
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

	results, errs := gotask.WhenAll(jobs, 0)

	for key, ret := range results {
		fmt.Println("key " + key + " has result of: " + ret.(string))
	}

	for key, err := range errs {
		fmt.Println("key " + key + " has error of: " + err.(error).Error())
	}
}
