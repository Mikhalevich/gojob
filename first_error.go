package jober

type FirstError struct {
	job
}

func NewFirstError() *FirstError {
	return &FirstError{
		job: *newJob(),
	}
}

func (self *FirstError) processError() {
	err, ok := <-self.errorChan
	if ok {
		self.dataErrors = append(self.dataErrors, err)
	}
	self.errorFinishFlag <- true
}

func (self *FirstError) processFunc() {
	self.waitGroup.Wait()
	close(self.dataChan)
	close(self.errorChan)
	<-self.dataFinishFlag
}

func (self *FirstError) Add(f WorkerFunc) {
	if self.startProcess(self) {
		go self.processFunc()
	}
	self.job.Add(f)
}

func (self *FirstError) addCallback(f WorkerFunc, callback func()) {
	if self.startProcess(self) {
		go self.processFunc()
	}
	self.job.addCallback(f, callback)
}

func (self *FirstError) Wait() {
	<-self.errorFinishFlag
}
