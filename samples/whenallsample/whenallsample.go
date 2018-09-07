package main

import (
	"fmt"

	"github.com/silver-xu/gotask"
)

func main() {

	jobs := []func() (interface{}, error){
		func() (interface{}, error) {
			return 1, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}

	results, _ := gotask.WhenAll(jobs, 1)

	for i := range results {
		fmt.Println(i)
	}
}
