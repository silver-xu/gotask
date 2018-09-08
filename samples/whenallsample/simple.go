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
