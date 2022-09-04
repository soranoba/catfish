package main

import (
	"github.com/soranoba/catfish-server/pkg/config"
	"github.com/soranoba/catfish-server/pkg/evaler"
	"github.com/soranoba/catfish-server/pkg/template"
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

func NewResponsePreset(name string, res *config.Response) (*ResponsePreset, error) {
	tpl, err := template.New(res.Body)
	if err != nil {
		return nil, err
	}

	cond := res.Condition
	if cond == "" {
		cond = "1.0"
	}

	return &ResponsePreset{
		Name:         name,
		Condition:    cond,
		Delay:        time.Duration(res.Delay * 1000_000_000),
		Status:       res.Status,
		Header:       res.Header,
		BodyTemplate: tpl,
	}, nil
}

func ElectResponsePreset(presets []*ResponsePreset, defaultPreset *ResponsePreset) *ResponsePreset {
	amountScore := float64(0)
	val := rand.Float64()

	evaler := evaler.New()
	for _, preset := range presets {
		score, _ := evaler.Eval(preset.Condition)
		if score+amountScore > val {
			return preset
		}
		amountScore += score
	}
	return defaultPreset
}
