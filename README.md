# gotask 

The Golang provides a easy and raw way to deal with concurrency with the legendary goroutine and channels. However from time to time it can be a bit "too raw" to spawn up extra goroutines, complete I/O execution then signal the channel and come back.

I have written this simple gotask package to allow C# coroutine style statement to automate the above process and simplify development.

It also allows optional cancellation of the task by passing boolean to the cancelChannel.

## Execute the following in terminal to install gotask

```
go get github.com/silver-xu/gotask
```

## Simple Usage:

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/silver-xu/gotask"
)

func main() {
	response, err := gotask.Await(doWork, nil)

	if err == nil {
		fmt.Println(response)
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

## Task Cancellation:

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/silver-xu/gotask"
)

func main() {
	cancelChannel := make(chan bool)

	response, err := gotask.Await(doWork, cancelChannel)

	cancelChannel <- true
}

func doWork() (interface{}, error) {
	//Todo: works
}

```