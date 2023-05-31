package configmap

import (
	"log"
)

var (
	errorf func(string, ...any) = log.Printf
	warnf  func(string, ...any) = log.Printf
	debugf func(string, ...any) = log.Printf
)

func SetErrorf(f func(string, ...any)) {
	if f != nil {
		errorf = f
	}
}

func SetWarnf(f func(string, ...any)) {
	if f != nil {
		warnf = f
	}
}

func SetDebugf(f func(string, ...any)) {
	if f != nil {
		debugf = f
	}
}
