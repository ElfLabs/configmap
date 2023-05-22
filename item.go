package configmap

import (
	"bytes"
	"encoding/json"
)

type Item struct {
}

type DecodeItemFunc func(in, out any, tag string) (changed bool, err error)

type nopLock struct{}

func (l nopLock) Lock() {
	return
}

func (l nopLock) Unlock() {
	return
}

func JsonDecodeItem(in, out any, _ string) (bool, error) {
	inBytes, err := json.Marshal(in)
	if err != nil {
		return false, err
	}
	outBytes, err := json.Marshal(out)
	if err != nil {
		return false, err
	}
	if bytes.Equal(inBytes, outBytes) {
		return false, nil
	}
	if err = json.Unmarshal(inBytes, out); err != nil {
		return false, err
	}
	return true, nil
}
