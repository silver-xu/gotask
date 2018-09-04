package gotask

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
