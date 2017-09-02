package jober

type WorkerFunc func() (interface{}, error)

type Jober interface {
	Add(f WorkerFunc)
	Wait()
}
