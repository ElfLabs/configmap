package configmap

var (
	configmap = New()
)

func Register(key string, item any, path ...string) *ConfigMap {
	return configmap.Register(key, item, path...)
}

func Load(provider Provider, parser Parser, opts ...Option) error {
	return configmap.Load(provider, parser, opts...)
}

func Get(path ...string) (any, bool) {
	return configmap.Get(path...)
}

func Display() string {
	return configmap.Display()
}
