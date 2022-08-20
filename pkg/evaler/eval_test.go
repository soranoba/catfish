package evaler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaler_Eval(t *testing.T) {
	assert := assert.New(t)

	val, err := New().Eval("0.8")
	assert.NoError(err)
	assert.Equal(0.8, val)

	val, err = New().Eval("0")
	assert.NoError(err)
	assert.Equal(0.0, val)

	val, err = New().Eval("1")
	assert.NoError(err)
	assert.Equal(1.0, val)

	val, err = New().Eval("true")
	assert.Error(err)
}
