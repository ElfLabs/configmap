package util

import (
	"sync"
)

type nopLock struct{}

func (l nopLock) Lock() {
	return
}

func (l nopLock) Unlock() {
	return
}

func NewNopLock() sync.Locker {
	return &nopLock{}
}
