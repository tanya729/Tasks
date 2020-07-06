package hw8

//WorkerFunction keep functions to workers
type WorkerFunction func() error

// Run will run slice of functions on workersCount workers while get maxErrors errors
func Run(functions []WorkerFunction, workersCount int, maxErrors int) ([]error, int) {
	tasksChan := make(chan WorkerFunction)
	resultsChan := make(chan error, workersCount-1)
	closeChan := make(chan bool)
	exitChan := make(chan bool)

	for i := 0; i < workersCount; i++ {
		go startWorker(tasksChan, resultsChan, closeChan, exitChan)
	}

	counter, errorsSlice := func() (int, []error) {
		var errors int
		var inProgress int
		var counter int
		var errorsSlice []error

		for i := 0; i < workersCount; i++ {
			inProgress++
			tasksChan <- functions[i]
		}

		for {
			err := <-resultsChan
			inProgress--
			counter++

			if err != nil {
				errorsSlice = append(errorsSlice, err)
				errors++
			}

			if counter == len(functions) || errors == maxErrors {
				close(closeChan)
				return counter, errorsSlice
			} else if len(functions)-counter-inProgress > 0 {
				inProgress++
				tasksChan <- functions[workersCount-1+counter]
			}
		}
	}()

	for i := 0; i < workersCount; i++ {
		<-exitChan
	}
	return errorsSlice, counter
}

func startWorker(tasksChan <-chan WorkerFunction, resultsChan chan<- error, closeChan <-chan bool, exitChan chan<- bool) {
	defer func() {
		exitChan <- true
	}()

	for {
		select {
		case task := <-tasksChan:
			resultsChan <- task()
		case <-closeChan:
			return
		}
	}
}
