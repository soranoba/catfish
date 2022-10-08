package main

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLogFormatter(t *testing.T) {
	var _ logrus.Formatter = &LogFormatter{}
}
