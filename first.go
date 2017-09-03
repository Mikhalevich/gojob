package jober

type First struct {
	Processor
}

func NewFirst() *First {
	return &First{
		Processor: *newProcessor(),
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

func (self *First) Add(workerFunc WorkerFunc) {
	if self.startProcess() {
		go self.processFunc()
	}

	self.waitGroup.Add(1)
	go func() {
		defer self.waitGroup.Done()
		d, err := workerFunc()
		if err != nil {
			self.errorChan <- err
			return
		}
		self.dataChan <- d
	}()
}

func (self *First) Wait() {
	<-self.dataFinishFlag
}
