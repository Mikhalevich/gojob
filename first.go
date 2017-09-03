package jober

type First struct {
	Processor
	funcChan chan WorkerFunc
	doneChan chan bool
}

func NewFirst() *First {
	return &First{
		Processor: *newProcessor(),
		funcChan:  make(chan WorkerFunc),
		doneChan:  make(chan bool),
	}
}

func (self *First) processFunc() {
	alive := true
	for {
		select {
		case f := <-self.funcChan:
			if alive {
				self.waitGroup.Add(1)
				go func() {
					defer self.waitGroup.Done()
					d, err := f()
					if err != nil {
						self.errorChan <- err
						return
					}
					self.dataChan <- d
				}()
			}
		case <-self.doneChan:
			alive = false
		}
	}
}

func (self *First) processData() {
	d, ok := <-self.dataChan
	if ok {
		self.data = append(self.data, d)
	}
	self.doneChan <- true
	self.dataFinishFlag <- true
}

func (self *First) Add(workerFunc WorkerFunc) {
	if self.startProcess() {
		go self.processFunc()
	}
	self.funcChan <- workerFunc
}
