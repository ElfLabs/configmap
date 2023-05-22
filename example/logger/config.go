package logger

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
