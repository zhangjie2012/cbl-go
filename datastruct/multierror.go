package datastruct

import (
	"errors"
	"strings"
	"sync"
)

// MultiError multiple errors, its useful in single error can not represent complex error
// all it's methods is gorountine safe.
// ref https://github.com/prometheus/alertmanager/blob/master/types/types.go#L268
type MultiError struct {
	mtx    sync.Mutex
	errors []error
}

func (e *MultiError) Add(err error) {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.errors = append(e.errors, err)
}

// AddStrError add string error
func (e *MultiError) AddStrError(err string) {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.errors = append(e.errors, errors.New(err))
}

func (e *MultiError) Len() int {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	return len(e.errors)
}

// Errors get a errors copy
func (e *MultiError) Errors() []error {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	return append(make([]error, 0, len(e.errors)), e.errors...)
}

func (e *MultiError) String(delimiter string) string {
	if delimiter == "" {
		delimiter = "; "
	}

	e.mtx.Lock()
	defer e.mtx.Unlock()

	es := make([]string, 0, len(e.errors))
	for _, err := range e.errors {
		es = append(es, err.Error())
	}
	return strings.Join(es, delimiter)
}
