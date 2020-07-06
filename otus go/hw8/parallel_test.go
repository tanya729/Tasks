package hw8

import (
	"errors"
	"testing"
	"time"
)

func TestRunWithoutErrors(t *testing.T) {
	f1 := func() error {
		time.Sleep(1 * time.Second)
		return nil
	}

	f2 := func() error {
		time.Sleep(2 * time.Second)
		return nil
	}

	f3 := func() error {
		time.Sleep(1 * time.Second)
		return nil
	}

	f4 := func() error {
		time.Sleep(2 * time.Second)
		return nil
	}

	f5 := func() error {
		time.Sleep(2 * time.Second)
		return nil
	}
	err, counter := Run([]WorkerFunction{
		f1,
		f2,
		f3,
		f4,
		f5,
	}, 2, 1)
	if len(err) > 0 {
		t.Errorf("Expected %d errors, got %d", 0, len(err))
	}
	if counter != 5 {
		t.Errorf("Expected %d tasks, got %d", 5, counter)
	}
}

func TestRunWithErrors(t *testing.T) {
	f1 := func() error {
		time.Sleep(1 * time.Second)
		return nil
	}

	f2 := func() error {
		time.Sleep(2 * time.Second)
		return errors.New("Some Error")
	}

	f3 := func() error {
		time.Sleep(1 * time.Second)
		return nil
	}

	f4 := func() error {
		time.Sleep(1 * time.Second)
		return nil
	}

	f5 := func() error {
		time.Sleep(2 * time.Second)
		return nil
	}

	workersCount := 2
	maxErrors := 1

	err, counter := Run([]WorkerFunction{
		f1,
		f2,
		f3,
		f4,
		f5,
	}, workersCount, maxErrors)
	if len(err) != 1 {
		t.Errorf("Expected %d errors, got %d", 1, len(err))
	}
	if counter > workersCount+maxErrors {
		t.Errorf("Expected %d tasks, got %d", workersCount+maxErrors, counter)
	}
}
