package jober

type First struct {
	job
}

func NewFirst() *First {
	return &First{
		job: *newJob(),
	}
}

func (f *First) processData() {
	d, ok := <-f.dataChan
	if ok {
		f.data = append(f.data, d)
		f.cancel()
	}
	f.dataFinishFlag <- true
}

func (f *First) processFunc() {
	f.waitGroup.Wait()
	close(f.dataChan)
	close(f.errorChan)
	<-f.errorFinishFlag
}

func (f *First) Add(fn WorkerFunc) {
	if f.startProcess(f) {
		go f.processFunc()
	}
	f.job.Add(fn)
}

func (f *First) addCallback(fn WorkerFunc, callback func()) {
	if f.startProcess(f) {
		go f.processFunc()
	}
	f.job.addCallback(fn, callback)
}

func (f *First) Wait() {
	<-f.dataFinishFlag
}
