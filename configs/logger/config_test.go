package logger

import (
	"testing"
)

func TestDevelopmentLogger(t *testing.T) {
	log, err := NewDevelopmentLogger().NewZapLogger()
	if err != nil {
		t.Errorf("new zap logger error: %s", err)
		return
	}
	log.Info("debug message")
	log.Info("info message")
}

func TestProductionLogger(t *testing.T) {
	log, err := NewProductionLogger().NewZapLogger()
	if err != nil {
		t.Errorf("new zap logger error: %s", err)
		return
	}
	log.Debug("should not display")
	log.Info("info message")
}
