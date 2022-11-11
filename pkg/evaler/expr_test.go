package evaler

import (
	"github.com/soranoba/valis"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestExpr_Eval(t *testing.T) {
	assert := assert.New(t)

	val, err := MustCompile("0.8").Eval(Params{})
	assert.NoError(err)
	assert.Equal(0.8, val)

	val, err = MustCompile("0").Eval(Params{})
	assert.NoError(err)
	assert.Equal(0.0, val)

	val, err = MustCompile("1").Eval(Params{})
	assert.NoError(err)
	assert.Equal(1.0, val)

	val, err = MustCompile("true").Eval(Params{})
	assert.NoError(err)
	assert.Equal(1.0, val)
}

func TestExpr_Eval_variables(t *testing.T) {
	assert := assert.New(t)

	val, err := MustCompile("x > 2").Eval(Params{"x": 1})
	assert.NoError(err)
	assert.Equal(0.0, val)

	val, err = MustCompile("x > 2").Eval(Params{"x": 3})
	assert.NoError(err)
	assert.Equal(1.0, val)

	val, err = MustCompile("x * 2").Eval(Params{"x": 0.1})
	assert.NoError(err)
	assert.Equal(0.2, val)

	val, err = MustCompile("x * 2").Eval(Params{"x": 1})
	assert.NoError(err)
	assert.Equal(2.0, val)
}

func TestExpr_Eval_queries(t *testing.T) {
	assert := assert.New(t)

	u, err := url.Parse("https://example.com?key=value&key=value2")
	assert.NoError(err)

	val, err := MustCompile("queries.key[0] == \"value\"").Eval(Params{"queries": u.Query()})
	assert.NoError(err)
	assert.Equal(1.0, val)
}

func TestExpr_Eval_atoi(t *testing.T) {
	assert := assert.New(t)
	val, err := MustCompile("atoi(v) == 123").Eval(Params{"v": "123"})
	assert.NoError(err)
	assert.Equal(1.0, val)
}

func TestExpr_UnmarshalText(t *testing.T) {
	assert := assert.New(t)
	expr := Expr{}
	assert.NoError(expr.UnmarshalText([]byte(".")))
	assert.Error(expr.Validate())
	assert.NoError(expr.UnmarshalText([]byte("x + 1")))
	assert.NoError(expr.Validate())

	val, err := expr.Eval(Params{"x": 2})
	assert.Equal(3.0, val)
	assert.NoError(err)
}

func TestExpr_MarshalText(t *testing.T) {
	assert := assert.New(t)
	expr := MustCompile("x + 1")
	val, err := expr.MarshalText()
	assert.NoError(err)
	assert.Equal("x + 1", string(val))
}

func TestExpr_Validate(t *testing.T) {
	var _ valis.Validatable = &Expr{}

	assert := assert.New(t)

	var expr Expr
	assert.NoError(expr.UnmarshalText([]byte("x x x")))
	assert.EqualError(
		expr.Validate(),
		"bad expression: 'x x x'",
	)
}
