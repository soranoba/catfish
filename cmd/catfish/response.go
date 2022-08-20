package main

import (
	"github.com/soranoba/catfish-server/pkg/config"
	"github.com/soranoba/catfish-server/pkg/evaler"
	"github.com/soranoba/catfish-server/pkg/template"
	"math/rand"
)

type (
	ResponsePreset struct {
		Name         string
		Condition    string
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

	return &ResponsePreset{
		Name:         name,
		Condition:    res.Condition,
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
