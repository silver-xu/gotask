package gotask

import (
	"errors"
	"sync"
)

//Await function equivalent to .Net Awaiter
func Await(do func() (interface{}, error), cancel chan bool) (interface{}, error) {
	taskChannel := make(chan interface{})
	errorChannel := make(chan error)

	go func(do func() (interface{}, error)) {
		select {
		case <-cancel:
			return
		default:
			result, err := do()

			taskChannel <- result
			errorChannel <- err

			close(taskChannel)
			close(errorChannel)

			return
		}
	}(do)

	return <-taskChannel, <-errorChannel
}

func WhenAll(doSlice []func() (interface{}, error), numberOfWorkers int) ([]interface{}, []error) {
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

	jobsChannel := make(chan func() (interface{}, error), len(doSlice))

	for i := 0; i < len(doSlice); i++ {
		jobsChannel <- doSlice[i]
	}

	close(jobsChannel)

	var worksWaitGroup sync.WaitGroup
	resultsChannel := make(chan interface{}, len(doSlice))
	errorsChannel := make(chan error, len(doSlice))

	worksWaitGroup.Add(len(doSlice))

	for i := 0; i < numberOfWorkers; i++ {

		//map
		go func(jobsChannel chan func() (interface{}, error), resultsChannel chan interface{}, errorsChannel chan error, wg *sync.WaitGroup) {

			for job := range jobsChannel {
				ret, err := job()

				resultsChannel <- ret
				errorsChannel <- err

				worksWaitGroup.Done()
			}

		}(jobsChannel, resultsChannel, errorsChannel, &worksWaitGroup)
	}

	//All works completed
	worksWaitGroup.Wait()

	//close all communication channels
	close(resultsChannel)
	close(errorsChannel)

	//reduce
	for result := range resultsChannel {
		results = append(results, result)
	}

	return results, errs
}
