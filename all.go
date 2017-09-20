package jober

type All struct {
	job
}

func NewAll() *All {
	return &All{
		job: *newJob(),
	}
}

func (self *All) Add(f WorkerFunc) {
	self.startProcess(self)
	self.job.Add(f)
}

func (self *All) addCallback(f WorkerFunc, callback func()) {
	self.startProcess(self)
	self.job.addCallback(f, callback)
}
