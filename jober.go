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

type processer interface {
	processData()
	processError()
}

type job struct {
	waitGroup       sync.WaitGroup
	active          bool
	data            []interface{}
	dataChan        chan interface{}
	dataFinishFlag  chan bool
	dataErrors      []error
	errorChan       chan error
	errorFinishFlag chan bool
	cancelationChan chan bool
}

func newJob() *job {
	return &job{
		active:          false,
		data:            make([]interface{}, 0),
		dataChan:        make(chan interface{}),
		dataFinishFlag:  make(chan bool),
		dataErrors:      make([]error, 0),
		errorChan:       make(chan error),
		errorFinishFlag: make(chan bool),
		cancelationChan: make(chan bool),
	}
}

func (j *job) cancel() {
	close(j.cancelationChan)
}

func (self *job) processData() {
	for d := range self.dataChan {
		self.data = append(self.data, d)
	}
	self.dataFinishFlag <- true
}

func (self *job) processError() {
	for err := range self.errorChan {
		self.dataErrors = append(self.dataErrors, err)
	}
	self.errorFinishFlag <- true
}

func (self *job) startProcess(p processer) bool {
	if !self.active {
		go p.processData()
		go p.processError()
		self.active = true
		return true
	}
	return false
}

func (self *job) Wait() {
	self.waitGroup.Wait()
	close(self.dataChan)
	close(self.errorChan)
	<-self.dataFinishFlag
	<-self.errorFinishFlag
}

func (self *job) Get() ([]interface{}, []error) {
	return self.data, self.dataErrors
}

func (self *job) Add(f WorkerFunc) {
	self.waitGroup.Add(1)
	go func() {
		defer self.waitGroup.Done()
		d, err := f()
		if err != nil {
			select {
			case self.errorChan <- err:
			case <-self.cancelationChan:
			}
			return
		}
		select {
		case self.dataChan <- d:
		case <-self.cancelationChan:
		}
	}()
}

func (self *job) addCallback(f WorkerFunc, callback func()) {
	self.waitGroup.Add(1)
	go func() {
		defer callback()
		defer self.waitGroup.Done()
		d, err := f()
		if err != nil {
			self.errorChan <- err
			return
		}
		self.dataChan <- d
	}()
}
