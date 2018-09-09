package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/silver-xu/gotask"
)

func main() {

	response, err := gotask.Await(doWork)

	if err == nil {
		fmt.Println(response)
	} else {
		fmt.Println(err)
	}
}

func doWork() (interface{}, error) {
	time.Sleep(time.Second * 5)
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
