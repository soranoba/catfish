package validator

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/tagrule"
	"github.com/soranoba/valis/when"
)

type (
	Validator struct {
		valis *valis.Validator
	}
)

func NewValidator() *Validator {
	v := valis.NewValidator()
	v.SetCommonRules(
		when.IsStruct(valis.EachFields(tagrule.Required, tagrule.Validate, tagrule.Enums)).
			ElseWhen(when.IsSliceOrArray(valis.Each( /* only common rules */ ))).
			ElseWhen(when.IsMap(valis.EachValues( /* only common rules */ ))),
	)

	return &Validator{
		valis: v,
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.valis.Validate(i)
}
