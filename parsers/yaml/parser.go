package yaml

import (
	"gopkg.in/yaml.v3"
)

// YAML implements a YAML parser.
type YAML struct{}

// Parser returns a YAML Parser.
func Parser() *YAML {
	return &YAML{}
}

func (p *YAML) Name() string {
	return "yaml"
}

// Unmarshal parses the given YAML bytes.
func (p *YAML) Unmarshal(b []byte, v any) error {
	return yaml.Unmarshal(b, v)
}

// Marshal marshals the given config map to YAML bytes.
func (p *YAML) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}
