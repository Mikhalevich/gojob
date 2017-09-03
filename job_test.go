package jober

import (
	"errors"
	"testing"
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

func TestFirstSimple(t *testing.T) {
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
