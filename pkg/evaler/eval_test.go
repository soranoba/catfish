package evaler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaler_Eval(t *testing.T) {
	assert := assert.New(t)

	val, err := New().Eval("0.8", Args{})
	assert.NoError(err)
	assert.Equal(0.8, val)

	val, err = New().Eval("0", Args{})
	assert.NoError(err)
	assert.Equal(0.0, val)

	val, err = New().Eval("1", Args{})
	assert.NoError(err)
	assert.Equal(1.0, val)

	val, err = New().Eval("true", Args{})
	assert.NoError(err)
	assert.Equal(1.0, val)
}

func TestEvaler_Eval_variables(t *testing.T) {
	assert := assert.New(t)

	val, err := New().Eval("x > 2", Args{"x": 1})
	assert.NoError(err)
	assert.Equal(0.0, val)

	val, err = New().Eval("x > 2", Args{"x": 3})
	assert.NoError(err)
	assert.Equal(1.0, val)

	val, err = New().Eval("x * 2", Args{"x": 0.1})
	assert.NoError(err)
	assert.Equal(0.2, val)

	val, err = New().Eval("x * 2", Args{"x": 1})
	assert.NoError(err)
	assert.Equal(2.0, val)
}
