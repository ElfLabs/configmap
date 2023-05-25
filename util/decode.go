package util

import (
	"bytes"
	"encoding/json"
)

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
