package main

import (
	"github.com/soranoba/catfish-server/pkg/config"
	"github.com/soranoba/henge/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElectResponsePreset(t *testing.T) {
	assert := assert.New(t)

	var responses [5]*config.Response
	for i, _ := range responses {
		responses[i] = &config.Response{
			Condition: "0.2",
			Status:    200 + i,
		}
	}

	var presets []*ResponsePreset
	for i, res := range responses {
		preset, err := NewResponsePreset(henge.ToString(i), res)
		if assert.NoError(err) {
			presets = append(presets, preset)
		}
	}

	totalElect := 100
	counts := make(map[string]int)
	for i := 0; i < totalElect; i++ {
		preset := ElectResponsePreset(presets[0:len(presets)-1], presets[len(presets)-1])
		if assert.NotNil(preset) {
			counts[preset.Name] += 1
		}
	}

	min := totalElect/len(presets) - 10
	max := totalElect/len(presets) + 10
	for _, preset := range presets {
		assert.True(min <= counts[preset.Name])
		assert.True(max >= counts[preset.Name])
	}
}
