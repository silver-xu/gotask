# gotask 

The Golang provides a easy and raw way to deal with concurrency with the legendary goroutine and channels. However from time to time it can be a bit "too raw" to spawn up extra goroutines, complete I/O execution then signal the channel and come back.

I have written this simple gotask package to allow C# coroutine style statement to automate the above process and simplify development.

Currently the project supports Await and WaitAll

## Await
Use c# style await to spawn a goroutine, finish the task and comeback. The main routine will wait for the result

It also allows optional cancellation of the task by passing boolean to the cancelChannel.

## WaitAll
Use C# style wait all to spawn multiple goroutine to finish a batch of jobs. The number of goroutine spawned can be controlled by numberOfWorkers flag.

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

## Await Task Cancellation:

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

## Simple WaitAll

```golang
package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/silver-xu/gotask"
)

func main() {
	jobs := map[string]func() (interface{}, error){
		"abc": func() (interface{}, error) {
			return 1, nil
		},
		"def": func() (interface{}, error) {
			return nil, errors.New("My Error")
		},
	}

	results, errs := gotask.WhenAll(jobs, 2)

	for key, ret := range results {
		fmt.Println("key " + key + " has result of: " + strconv.Itoa(ret.(int)))
	}

	for key, err := range errs {
		fmt.Println("key " + key + " has error of: " + err.(error).Error())
	}
}
```
