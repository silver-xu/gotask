package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/silver-xu/gotask"
)

func main() {

	url := "https://www.google.com/"

	cancelChannel := make(chan bool)

	response, err := gotask.Await(func() (interface{}, error) {
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

	}, cancelChannel)

	if err == nil {
		fmt.Println(response)
	}
}
