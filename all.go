package jober

import (
	"sync"
)

type All struct {
	waitGroup sync.WaitGroup
	active    bool
	dataFlag  chan bool
	dataChan  chan interface{}
	errorFlag chan bool
	errorChan chan error
	Errors    []error
	Results   []interface{}
}

func NewAll() *All {
	return &All{
		active:    false,
		dataFlag:  make(chan bool),
		dataChan:  make(chan interface{}),
		errorFlag: make(chan bool),
		errorChan: make(chan error),
		Results:   make([]interface{}, 0),
		Errors:    make([]error, 0),
	}
}

func (self *All) processData() {
	for data := range self.dataChan {
		self.Results = append(self.Results, data)
	}
	self.dataFlag <- true
}

func (self *All) processError() {
	for error := range self.errorChan {
		self.Errors = append(self.Errors, error)
	}
	self.errorFlag <- true
}

func (self *All) Add(workerFunc WorkerFunc) {
	if !self.active {
		go self.processData()
		go self.processError()
		self.active = true
	}

	self.waitGroup.Add(1)
	go func() {
		defer self.waitGroup.Done()
		data, err := workerFunc()
		if err != nil {
			self.errorChan <- err
			return
		}
		self.dataChan <- data
	}()
}

func (self *All) Wait() {
	self.waitGroup.Wait()
	close(self.dataChan)
	close(self.errorChan)
	<-self.dataFlag
	<-self.errorFlag
}
