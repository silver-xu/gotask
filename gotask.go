package gotask

import (
	"context"
	"errors"
	"sync"
	"time"
)

type kvPair struct {
	key   string
	value interface{}
}

//Await function equivalent to .Net Awaiter
func Await(do func() (interface{}, error), timeout time.Duration) (interface{}, error) {
	resultChannel := make(chan interface{})
	errorChannel := make(chan error)

	go func(do func() (interface{}, error)) {

		result, err := do()
		resultChannel <- result
		errorChannel <- err

		close(resultChannel)
		close(errorChannel)

		return

	}(do)

	select {
	case result := <-resultChannel:
		return result, nil
	case err := <-errorChannel:
		return nil, err
	case <-time.After(time.Second * timeout):
		return nil, errors.New("Timeout Exception")
	}
}

/*WhenAll does: Assign the tasks in jobsSlice to multiple workers throttled by numberOfWorkers parameter
with a timeout value in seconds*/
func WhenAll(jobsSlice map[string]func() (interface{}, error), numberOfWorkers int, timeout time.Duration) (map[string]interface{}, map[string]error) {

	results := make(map[string]interface{})
	errs := make(map[string]error)

	//parameters validation and correction
	if jobsSlice == nil {
		errs[""] = errors.New("Argument Exception")

		return results, errs
	}

	//set numberOfWorkers no to exceed the jobs to be done
	if numberOfWorkers == 0 || numberOfWorkers > len(jobsSlice) {
		numberOfWorkers = len(jobsSlice)
	}

	//make the job queue
	jobsChannel := make(chan kvPair, len(jobsSlice))

	for key, function := range jobsSlice {
		jobsChannel <- kvPair{key, function}
	}

	close(jobsChannel)

	var workersWaitGroup sync.WaitGroup

	resultsChannel := make(chan kvPair, len(jobsSlice))
	errorsChannel := make(chan kvPair, len(jobsSlice))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	defer cancel()

	workersWaitGroup.Add(len(jobsSlice))

	for i := 0; i < numberOfWorkers; i++ {

		//map
		go func(jobsChannel chan kvPair, resultsChannel chan kvPair, errorsChannel chan kvPair, ctx context.Context, wg *sync.WaitGroup) {
			for jobPair := range jobsChannel {
				work := jobPair.value.(func() (interface{}, error))
				ret, err := work()

				select {
				case <-time.After(1 * time.Nanosecond):
					if err != nil {
						errorsChannel <- kvPair{jobPair.key, err}
					} else {
						resultsChannel <- kvPair{jobPair.key, ret}
					}

					workersWaitGroup.Done()

				// we received the signal of cancelation in this channel
				case <-ctx.Done():
					workersWaitGroup.Done()

					errorsChannel <- kvPair{"", errors.New("Timeout Exception")}
				}

			}

		}(jobsChannel, resultsChannel, errorsChannel, ctx, &workersWaitGroup)
	}

	//All works completed
	workersWaitGroup.Wait()

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
