package jober

type All struct {
	Processor
}

func NewAll() *All {
	return &All{
		Processor: *newProcessor(),
	}
}

func (self *All) Add(workerFunc WorkerFunc) {
	self.startProcess()
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
