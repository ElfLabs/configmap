package logger

import (
	"fmt"
	"log"
	"reflect"
)

type Logger struct {
	Level      string            `json:"level,omitempty" yaml:"level,omitempty"`
	Console    bool              `json:"console" yaml:"console"`
	Filename   string            `json:"filename,omitempty" yaml:"filename,omitempty"`
	InitFields map[string]string `json:"initFields,omitempty" yaml:"initFields,omitempty"`
	Outputs    []string          `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

func NewDefaultConfig() *Logger {
	return &Logger{
		Level:    "info",
		Console:  true,
		Filename: "out.log",
		Outputs: []string{
			"stderr",
		},
	}
}

func (l *Logger) Update(v any) error {
	tmp, ok := v.(*Logger)
	if !ok {
		return fmt.Errorf("unexpect type: %s", reflect.TypeOf(l))
	}
	*l = *tmp
	log.Printf("logger updated!")
	return nil
}

func (l *Logger) Changed() {
	log.Printf("logger changed!")
}
