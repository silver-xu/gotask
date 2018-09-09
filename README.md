# gotask 

The Golang provides a easy and raw way to deal with concurrency with the legendary goroutine and channels. However from time to time it can be a bit "too raw" to spawn up extra goroutines, complete I/O execution then signal the channel and come back.

I have written this simple gotask package to allow C# coroutine style statement to automate the above process and simplify development.

Currently the project supports Await and WaitAll

## Await
Use c# style await to spawn a goroutine, finish the task and comeback. The main routine will wait for the result

It also allows optional timeout of the tasks after n Seconds, if timeout is not specified then it will never timeout

## WhenAll
Use C# style wait all to spawn multiple goroutine to finish a batch of jobs. The number of goroutine spawned can be controlled by numberOfWorkers flag.

It also allows optional timeout of the tasks after n Seconds, if timeout is not specified then it will never timeout

## Execute the following in terminal to install gotask

```
go get github.com/silver-xu/gotask
```

## Simple Await:

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/silver-xu/gotask"
)

func main() {
    //timeout after 1 seconds
	response, err := gotask.Await(doWork, 1)

	if err == nil {
		fmt.Println(response)
	} else {
		fmt.Println(err)
	}
}

func doWork() (interface{}, error) {
	url := "https://www.google.com/"

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
```

## WhenAll

```golang
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
```
