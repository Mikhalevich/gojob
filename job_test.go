package gojob

import (
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
	for _, r := range job.Results {
		testSumm += r.(int)
	}

	if summ != testSumm {
		t.Fatalf("Compared values are not the same %d -> %d", summ, testSumm)
	}
}