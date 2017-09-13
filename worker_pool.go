package jober

type WorkerPool struct {
	job   Jober
	count chan bool
}

func NewWorkerPool(j Jober, c int) *WorkerPool {
	return &WorkerPool{
		job:   j,
		count: make(chan bool, c),
	}
}

func (wp *WorkerPool) Add(f WorkerFunc) {
	wp.count <- true
	wp.job.Add(f)
}

func (wp *WorkerPool) Wait() {
	close(wp.count)
	wp.job.Wait()
}

func (wp *WorkerPool) Get() ([]interface{}, []error) {
	return wp.job.Get()
}
