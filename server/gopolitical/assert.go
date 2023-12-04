package gopolitical

import (
	"reflect"
	"testing"
)

type Assert struct {
	t *testing.T
}

func NewAssert(t *testing.T) Assert {
	return Assert{t}
}

func (a *Assert) NoError(err error) {
	a.t.Helper() // increase stack pointer in log
	if err != nil {
		a.t.Errorf("Unexpected error: %v", err)
	}
}

func (a *Assert) Error(err error) {
	a.t.Helper() // increase stack pointer in log
	if err == nil {
		a.t.Error("An error was expected")
	}
}

func (a *Assert) True(boolean bool) {
	a.t.Helper() // increase stack pointer in log
	if boolean {
		a.t.Error("Boolean must be true")
	}
}

func (a *Assert) DeepEqual(got any, expected any) {
	a.t.Helper() // increase stack pointer in log
	if !reflect.DeepEqual(got, expected) {
		a.t.Errorf("Results mismatch: Got %v, Expected %v", got, expected)
	}
}

func (a *Assert) Equal(got any, expected any) {
	a.t.Helper() // increase stack pointer in log
	if got != expected {
		a.t.Errorf("Results mismatch: Got %v, Expected %v", got, expected)
	}
}

func (a *Assert) Empty(got any) {
	a.t.Helper()   // increase stack pointer in log
	defer func() { // recover from failed len()
		if r := recover(); r != nil {
			a.t.Errorf("Expected empty array: Got %v", got)
		}
	}()
	if reflect.ValueOf(got).Len() != 0 {
		a.t.Errorf("Expected empty array: Got %v", got)
	}
}
