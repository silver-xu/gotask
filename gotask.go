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
		for {
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
		}
	}(do)

	return <-taskChannel, <-errorChannel
}

type kvPair struct {
	key   string
	value interface{}
}

/*WhenAll does: Assign the tasks in doSlice to multiple workers throttled by numberOfWorkers parameter */
func WhenAll(doSlice map[string]func() (interface{}, error), numberOfWorkers int) (map[string]interface{}, map[string]error) {
	results := make(map[string]interface{})
	errs := make(map[string]error)

	if doSlice == nil {
		errs[""] = errors.New("Argument Exception")

		return results, errs
	}

	//set numberOfWorkers no to exceed the jobs to be done
	if numberOfWorkers == 0 || numberOfWorkers > len(doSlice) {
		numberOfWorkers = len(doSlice)
	}

	jobsChannel := make(chan kvPair, len(doSlice))

	for key, function := range doSlice {
		jobsChannel <- kvPair{key, function}
	}

	close(jobsChannel)

	var worksWaitGroup sync.WaitGroup
	resultsChannel := make(chan kvPair, len(doSlice))
	errorsChannel := make(chan kvPair, len(doSlice))

	worksWaitGroup.Add(len(doSlice))

	for i := 0; i < numberOfWorkers; i++ {

		//map
		go func(jobsChannel chan kvPair, resultsChannel chan kvPair, errorsChannel chan kvPair, wg *sync.WaitGroup) {

			for jobPair := range jobsChannel {

				work := jobPair.value.(func() (interface{}, error))
				ret, err := work()

				if err != nil {
					errorsChannel <- kvPair{jobPair.key, err}
				} else {
					resultsChannel <- kvPair{jobPair.key, ret}
				}

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
	for resultKVPair := range resultsChannel {
		results[resultKVPair.key] = resultKVPair.value
	}

	for errKVPair := range errorsChannel {
		errs[errKVPair.key] = errKVPair.value.(error)
	}

	return results, errs
}
