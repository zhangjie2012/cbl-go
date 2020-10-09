package datastruct

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiError(t *testing.T) {
	m := MultiError{}
	assert.Equal(t, 0, m.Len())

	m.Add(errors.New("hello1"))
	m.AddStrError("hello2")
	assert.Equal(t, 2, m.Len())

	errs := m.Errors()
	assert.EqualValues(t, []error{errors.New("hello1"), errors.New("hello2")}, errs)

	errstr := m.String("")
	assert.Equal(t, "hello1; hello2", errstr)
}
