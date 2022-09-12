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
		Condition    string
		Delay        time.Duration
		Status       int
		Header       map[string]string
		BodyTemplate *template.Template
	}
)

func NewResponsePreset(res *config.Response) (*ResponsePreset, error) {
	tpl, err := template.New(res.Body)
	if err != nil {
		return nil, err
	}

	cond := res.Condition
	if cond == "" {
		cond = "1.0"
	}

	return &ResponsePreset{
		Name:         res.Name,
		Condition:    cond,
		Delay:        time.Duration(res.Delay * 1000_000_000),
		Status:       res.Status,
		Header:       res.Header,
		BodyTemplate: tpl,
	}, nil
}

func ElectResponsePreset(presets []*ResponsePreset, args evaler.Args) (*ResponsePreset, error) {
	amountScore := float64(0)
	val := rand.Float64()

	evaler := evaler.New()
	for _, preset := range presets {
		score, err := evaler.Eval(preset.Condition, args)
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
