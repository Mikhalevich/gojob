package gojob

import (
	"sync"
)

type First struct {
	waitGroup sync.WaitGroup
	active    bool
	dataFlag  chan bool
	dataChan  chan interface{}
	errorFlag chan bool
	errorChan chan error
	Errors    []error
	Results   interface{}
	funcChan  chan WorkerFunc
	doneChan  chan bool
}

func NewFirst() *First {
	return &First{
		active:    false,
		dataFlag:  make(chan bool),
		dataChan:  make(chan interface{}),
		errorFlag: make(chan bool),
		errorChan: make(chan error),
		Errors:    make([]error, 0),
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
					data, err := f()
					if err != nil {
						self.errorChan <- err
						return
					}
					self.dataChan <- data
				}()
			}
		case <-self.doneChan:
			alive = false
		}
	}
}

func (self *First) processData() {
	self.Results = <-self.dataChan
	self.doneChan <- true
	self.dataFlag <- true
}

func (self *First) processError() {
	for error := range self.errorChan {
		self.Errors = append(self.Errors, error)
	}
	self.errorFlag <- true
}

func (self *First) Add(workerFunc WorkerFunc) {
	if !self.active {
		go self.processFunc()
		go self.processData()
		go self.processError()
		self.active = true
	}

	self.funcChan <- workerFunc
}

func (self *First) Wait() {
	self.waitGroup.Wait()
	close(self.dataChan)
	close(self.errorChan)
	<-self.dataFlag
	<-self.errorFlag
}
