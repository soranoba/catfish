package evaler

import (
	"context"
	"errors"
	"fmt"
	"github.com/PaesslerAG/gval"
	"strconv"
)

type (
	Expr struct {
		raw        string
		expression gval.Evaluable
	}
	Params map[string]interface{}
)

func MustCompile(expr string) *Expr {
	evaler, err := Compile(expr)
	if err != nil {
		panic(err)
	}
	return evaler
}

func Compile(expr string) (*Expr, error) {
	expression, err := compile(expr)
	if err != nil {
		return nil, err
	}

	return &Expr{
		raw:        expr,
		expression: expression,
	}, nil
}

func (expr *Expr) Eval(args Params) (float64, error) {
	if expr.expression == nil {
		return 1.0, nil
	}

	value, err := expr.expression(context.Background(), args)
	if err != nil {
		return 0.0, err
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case bool:
		if v == true {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("return type of expr is invalid: %v", value)
	}
}

func (expr *Expr) UnmarshalText(s []byte) error {
	if len(s) == 0 {
		return nil
	}

	raw := string(s)
	expression, _ := compile(raw)
	*expr = Expr{
		expression: expression,
		raw:        raw,
	}
	return nil
}

func (expr Expr) MarshalText() ([]byte, error) {
	return []byte(expr.String()), nil
}

func (expr *Expr) Validate() error {
	if expr == nil {
		return nil
	}

	if _, err := compile(expr.raw); err != nil {
		return fmt.Errorf("bad expression: '%s'", expr.raw)
	}
	return nil
}

func (expr *Expr) String() string {
	return expr.raw
}

func compile(expr string) (gval.Evaluable, error) {
	return gval.Full(
		gval.Function("atoi", func(args ...interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, errors.New("expected exactly 1 argument")
			}
			if str, ok := args[0].(string); ok {
				val, err := strconv.Atoi(str)
				if err != nil {
					return nil, err
				}
				return float64(val), nil
			}
			return nil, errors.New("expected string")
		}),
	).NewEvaluable(expr)
}
