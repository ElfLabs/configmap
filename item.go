package configmap

type DecodeItemFunc func(in, out any, tag string) (changed bool, err error)

type Item interface {
	Key() string
}

type UpdateEvent interface {
	Update(any) error
}

type NotifyEvent interface {
	Changed()
}
