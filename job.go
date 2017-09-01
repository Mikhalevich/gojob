package gojob

type WorkerFunc func() (interface{}, error)

type Job interface {
	Add(f WorkerFunc)
	Wait()
}
