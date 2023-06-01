package toml

import (
	"github.com/pelletier/go-toml/v2"
)

// Toml implements a Toml parser.
type Toml struct{}

// Parser returns a Toml Parser.
func Parser() *Toml {
	return &Toml{}
}

func (p *Toml) Name() string {
	return "yaml"
}

// Unmarshal parses the given Toml bytes.
func (p *Toml) Unmarshal(b []byte, v any) error {
	return toml.Unmarshal(b, v)
}

// Marshal marshals the given config map to Toml bytes.
func (p *Toml) Marshal(v any) ([]byte, error) {
	return toml.Marshal(v)
}
