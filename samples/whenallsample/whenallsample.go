package main

import (
	"errors"
	"fmt"
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

	results, _ := whenAll(jobs, 0)

	for i := range results {
		fmt.Println(i)
	}

	fmt.Scanln()
}

func whenAll(doSlice []func() (interface{}, error), numberOfWorkers int) ([]interface{}, []error) {
	var results []interface{}
	var errs []error

	if doSlice == nil {
		errs = append(errs, errors.New("Argument Exception"))

		return results, errs
	}

	//set numberOfWorkers no to exceed the jobs to be done
	if numberOfWorkers == 0 || numberOfWorkers > len(doSlice) {
		numberOfWorkers = len(doSlice)
	}

	jobsChannel := make(chan func() (interface{}, error))

	for i := 0; i < numberOfWorkers; i++ {

		//map
		go func(jobsChannel chan func() (interface{}, error)) {

			for job := range jobsChannel {

				job()
			}

		}(jobsChannel)
	}

	for i := 0; i < len(doSlice); i++ {
		jobsChannel <- doSlice[i]
	}

	close(jobsChannel)

	return results, errs

}
