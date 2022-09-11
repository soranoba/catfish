package evaler

import (
	"errors"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"
)

type (
	Evaler struct {
	}
	Args map[string]interface{}
)

func New() *Evaler {
	return &Evaler{}
}

func (ev *Evaler) Eval(expr string, args Args) (float64, error) {
	pkg := types.NewPackage("main", "main")
	for name, value := range args {
		switch v := value.(type) {
		case int:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case int8:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case int16:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case int32:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case int64:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case uint:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case uint8:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case uint16:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case uint32:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case uint64:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case float32:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case float64:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Float64], constant.MakeFloat64(float64(v))))
		case bool:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.Bool], constant.MakeBool(v)))
		case string:
			pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, types.Typ[types.String], constant.MakeString(v)))
		default:
			break
		}
	}

	v, err := types.Eval(token.NewFileSet(), pkg, token.NoPos, expr)
	if err != nil {
		return 0.0, err
	}

	switch v.Value.Kind() {
	case constant.Int, constant.Float:
		return strconv.ParseFloat(v.Value.String(), 64)
	case constant.Bool:
		b, err := strconv.ParseBool(v.Value.String())
		if b {
			return 1.0, err
		}
		return 0.0, err
	default:
		return 0.0, errors.New("invalid expr")
	}
}
