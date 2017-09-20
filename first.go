package jober

type First struct {
	job
}

func NewFirst() *First {
	return &First{
		job: *newJob(),
	}
}

func (self *First) processData() {
	d, ok := <-self.dataChan
	if ok {
		self.data = append(self.data, d)
	}
	self.dataFinishFlag <- true
}

func (self *First) processFunc() {
	self.waitGroup.Wait()
	close(self.dataChan)
	close(self.errorChan)
	<-self.errorFinishFlag
}

func (self *First) Add(f WorkerFunc) {
	if self.startProcess(self) {
		go self.processFunc()
	}
	self.job.Add(f)
}

func (self *First) addCallback(f WorkerFunc, callback func()) {
	if self.startProcess(self) {
		go self.processFunc()
	}
	self.job.addCallback(f, callback)
}

func (self *First) Wait() {
	<-self.dataFinishFlag
}
