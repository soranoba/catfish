package evaler

import (
	"strconv"
)

type (
	Evaler struct {
	}
)

func New() *Evaler {
	return &Evaler{}
}

func (ev *Evaler) Eval(expr string) (float64, error) {
	return strconv.ParseFloat(expr, 64)
}
