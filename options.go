package configmap

import (
	"github.com/ElfLabs/configmap/util"
)

type Option func(opts *Options)
type Options struct {
	DecodeItemFunc DecodeItemFunc
	Providers      []Provider
}

func (o *Options) apply(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func newOptions(opts ...Option) Options {
	var o = Options{
		DecodeItemFunc: util.JsonDecodeItem,
	}
	o.apply(opts...)
	return o
}

func WithProvider(provider Provider) Option {
	return func(opts *Options) {
		opts.Providers = append(opts.Providers, provider)
	}
}
