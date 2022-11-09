package main

import (
	"github.com/soranoba/catfish/pkg/config"
	"github.com/soranoba/catfish/pkg/evaler"
	"github.com/soranoba/catfish/pkg/template"
	"math/rand"
	"time"
)

type (
	ResponsePreset struct {
		Name         string
		Condition    *evaler.Expr
		Delay        time.Duration
		Status       int
		Redirect     *string
		Header       map[string]string
		BodyTemplate *template.Template
	}
)

func NewResponsePreset(res *config.Response) (*ResponsePreset, error) {
	tpl, err := template.New(res.Body)
	if err != nil {
		return nil, err
	}

	return &ResponsePreset{
		Name:         res.Name,
		Condition:    res.Condition,
		Delay:        time.Duration(res.Delay * 1000_000_000),
		Status:       res.Status,
		Redirect:     res.Redirect,
		Header:       res.Header,
		BodyTemplate: tpl,
	}, nil
}

func ElectResponsePreset(presets []*ResponsePreset, args evaler.Params) (*ResponsePreset, error) {
	amountScore := float64(0)
	val := rand.Float64()

	for _, preset := range presets {
		// NOTE: Default always matches.
		if preset.Condition == nil {
			return preset, nil
		}

		score, err := preset.Condition.Eval(args)
		if err != nil {
			return nil, err
		}
		if score+amountScore > val {
			return preset, nil
		}
		amountScore += score
	}
	return nil, nil
}
