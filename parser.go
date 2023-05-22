package configmap

// Parser represents a configuration format parser.
type Parser interface {
	Name() string
	Unmarshal([]byte, any) error
	Marshal(any) ([]byte, error)
}

type parser struct {
	name      string
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

func (p parser) Name() string {
	return p.name
}

func (p parser) Unmarshal(b []byte, v any) error {
	return p.unmarshal(b, v)
}

func (p parser) Marshal(v any) ([]byte, error) {
	return p.marshal(v)
}

func WrapParser(name string, marshal func(any) ([]byte, error), unmarshal func([]byte, any) error) Parser {
	return &parser{
		name:      name,
		marshal:   marshal,
		unmarshal: unmarshal,
	}
}
