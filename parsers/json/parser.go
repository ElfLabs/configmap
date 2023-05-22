package json

import (
	"encoding/json"
)

// JSON implements a JSON parser.
type JSON struct{}

// Parser returns a JSON Parser.
func Parser() *JSON {
	return &JSON{}
}

func (p *JSON) Name() string {
	return "json"
}

// Unmarshal parses the given JSON bytes.
func (p *JSON) Unmarshal(b []byte, v any) error {
	return json.Unmarshal(b, v)
}

// Marshal marshals the given config map to JSON bytes.
func (p *JSON) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
