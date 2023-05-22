package configmap

// Provider represents a configuration provider. Providers can
// read configuration from a source (file, HTTP etc.)
type Provider interface {
	// ReadBytes returns the entire configuration as raw []bytes to be parsed.
	// with a Parser.
	ReadBytes() ([]byte, error)
}
