package jober

import (
	"sync"
)

type WorkerFunc func() (interface{}, error)

type Jober interface {
	Add(f WorkerFunc)
	Wait()
	Get() ([]interface{}, []error)
}

type Processor struct {
	waitGroup       sync.WaitGroup
	active          bool
	data            []interface{}
	dataChan        chan interface{}
	dataFinishFlag  chan bool
	dataErrors      []error
	errorChan       chan error
	errorFinishFlag chan bool
}

func newProcessor() *Processor {
	return &Processor{
		active:          false,
		data:            make([]interface{}, 0),
		dataChan:        make(chan interface{}),
		dataFinishFlag:  make(chan bool),
		dataErrors:      make([]error, 0),
		errorChan:       make(chan error),
		errorFinishFlag: make(chan bool),
	}
}

func (self *Processor) processData() {
	for d := range self.dataChan {
		self.data = append(self.data, d)
	}
	self.dataFinishFlag <- true
}

func (self *Processor) processError() {
	for err := range self.errorChan {
		self.dataErrors = append(self.dataErrors, err)
	}
	self.errorFinishFlag <- true
}

func (self *Processor) startProcess() bool {
	if !self.active {
		go self.processData()
		go self.processError()
		self.active = true
		return true
	}
	return false
}

func (self *Processor) Wait() {
	self.waitGroup.Wait()
	close(self.dataChan)
	close(self.errorChan)
	<-self.dataFinishFlag
	<-self.errorFinishFlag
}

func (self *Processor) Get() ([]interface{}, []error) {
	return self.data, self.dataErrors
}
