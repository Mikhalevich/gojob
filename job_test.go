package jober

import (
	"errors"
	"testing"
	"time"
)

func TestAllSumm(t *testing.T) {
	job := NewAll()
	summ := 0
	for i := 0; i < 100; i++ {
		summ += i
		locValue := i
		f := func() (interface{}, error) {
			return locValue, nil
		}
		job.Add(f)
	}

	job.Wait()

	testSumm := 0
	results, errors := job.Get()

	t.Logf("results len = %d, errors len = %d", len(results), len(errors))

	for _, r := range results {
		testSumm += r.(int)
	}

	if summ != testSumm {
		t.Fatalf("Compared values are not the same %d -> %d", summ, testSumm)
	}
}

func TestFirstStart(t *testing.T) {
	job := NewFirst()
	for i := 0; i < 100; i++ {
		locValue := i
		f := func() (interface{}, error) {
			if locValue == 0 {
				return 50, nil
			}
			time.Sleep(1000 * time.Millisecond)
			return 0, errors.New("Invalid value")
		}
		job.Add(f)
	}

	job.Wait()

	results, errors := job.Get()

	t.Logf("results len = %d, errors len = %d", len(results), len(errors))

	if len(results) != 1 {
		t.Fatalf("Not valid data length len = %d", len(results))
	}

	if results[0].(int) != 50 {
		t.Fatalf("Compared values are not the same %d -> %d", 50, results[0].(int))
	}
}

func TestFirstMiddle(t *testing.T) {
	job := NewFirst()
	for i := 0; i < 100; i++ {
		locValue := i
		f := func() (interface{}, error) {
			if locValue == 50 {
				return 50, nil
			}
			return 0, errors.New("Invalid value")
		}
		job.Add(f)
	}

	job.Wait()

	results, errors := job.Get()

	t.Logf("results len = %d, errors len = %d", len(results), len(errors))

	if len(results) != 1 {
		t.Fatalf("Not valid data length len = %d", len(results))
	}

	if results[0].(int) != 50 {
		t.Fatalf("Compared values are not the same %d -> %d", 50, results[0].(int))
	}
}

func TestFirstEnd(t *testing.T) {
	job := NewFirst()
	for i := 0; i < 100; i++ {
		locValue := i
		f := func() (interface{}, error) {
			if locValue == 99 {
				return 50, nil
			}
			return 0, errors.New("Invalid value")
		}
		job.Add(f)
	}

	job.Wait()

	results, errors := job.Get()

	t.Logf("results len = %d, errors len = %d", len(results), len(errors))

	if len(results) != 1 {
		t.Fatalf("Not valid data length len = %d", len(results))
	}

	if results[0].(int) != 50 {
		t.Fatalf("Compared values are not the same %d -> %d", 50, results[0].(int))
	}
}

func TestFirstErrorStart(t *testing.T) {
	job := NewFirstError()
	for i := 0; i < 100; i++ {
		locValue := i
		f := func() (interface{}, error) {
			if locValue == 0 {
				return 0, errors.New("First invalid value")
			}
			time.Sleep(1000 * time.Millisecond)
			return 50, nil
		}
		job.Add(f)
	}

	job.Wait()

	results, errors := job.Get()

	t.Logf("results len = %d, errors len = %d", len(results), len(errors))

	if len(errors) != 1 {
		t.Fatalf("Not valid data length len = %d", len(errors))
	}
}

func TestFirstErrorMiddle(t *testing.T) {
	job := NewFirstError()
	for i := 0; i < 100; i++ {
		locValue := i
		f := func() (interface{}, error) {
			if locValue == 50 {
				return 0, errors.New("First invalid value")
			}
			return 50, nil
		}
		job.Add(f)
	}

	job.Wait()

	results, errors := job.Get()

	t.Logf("results len = %d, errors len = %d", len(results), len(errors))

	if len(errors) != 1 {
		t.Fatalf("Not valid data length len = %d", len(errors))
	}
}

func TestFirstErrorEnd(t *testing.T) {
	job := NewFirstError()
	for i := 0; i < 100; i++ {
		locValue := i
		f := func() (interface{}, error) {
			if locValue == 99 {
				return 0, errors.New("First invalid value")
			}
			return 50, nil
		}
		job.Add(f)
	}

	job.Wait()

	results, errors := job.Get()

	t.Logf("results len = %d, errors len = %d", len(results), len(errors))

	if len(errors) != 1 {
		t.Fatalf("Not valid data length len = %d", len(errors))
	}
}
