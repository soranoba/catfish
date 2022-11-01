package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type (
	LogFormatter struct{}
	logEntryJson struct {
		Type    string                 `json:"type"`
		Level   string                 `json:"level"`
		Time    string                 `json:"time"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
)

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	jsonEntry := &logEntryJson{
		Type:    "default",
		Level:   entry.Level.String(),
		Time:    entry.Time.Format("2006-01-02T15:04:05.999Z07:00"),
		Message: entry.Message,
	}

	if logType, ok := entry.Data["@type"].(string); ok {
		jsonEntry.Type = logType
	}

	data := make(logrus.Fields, 0)
	for k, v := range entry.Data {
		if strings.HasPrefix(k, "@") {
			continue
		}
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	jsonEntry.Data = data

	b := entry.Buffer
	if b == nil {
		b = &bytes.Buffer{}
	}

	encoder := json.NewEncoder(b)
	if err := encoder.Encode(jsonEntry); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}
	return b.Bytes(), nil
}
