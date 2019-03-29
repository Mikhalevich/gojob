package jober

type FirstError struct {
	job
}

func NewFirstError() *FirstError {
	return &FirstError{
		job: *newJob(),
	}
}

func (fr *FirstError) processError() {
	err, ok := <-fr.errorChan
	if ok {
		fr.dataErrors = append(fr.dataErrors, err)
		fr.cancel()
	}
	fr.errorFinishFlag <- true
}

func (fr *FirstError) processFunc() {
	fr.waitGroup.Wait()
	close(fr.dataChan)
	close(fr.errorChan)
	<-fr.dataFinishFlag
}

func (fr *FirstError) Add(f WorkerFunc) {
	if fr.startProcess(fr) {
		go fr.processFunc()
	}
	fr.job.Add(f)
}

func (fr *FirstError) addCallback(f WorkerFunc, callback func()) {
	if fr.startProcess(fr) {
		go fr.processFunc()
	}
	fr.job.addCallback(f, callback)
}

func (fr *FirstError) Wait() {
	<-fr.errorFinishFlag
}
