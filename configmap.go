package configmap

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/armon/go-radix"
	"github.com/mitchellh/copystructure"

	"github.com/ElfLabs/configmap/util"
)

const (
	PathSeparator = "."
)

type ConfigMap struct {
	tree   *radix.Tree
	locker sync.RWMutex
}

func New() *ConfigMap {
	return &ConfigMap{
		tree: radix.New(),
	}
}

func getPath(key string, path ...string) string {
	switch len(path) {
	case 0:
		return key
	case 1:
		return path[0] + PathSeparator + key
	default:
		return strings.Join(append(path, key), PathSeparator)
	}
}

func (m *ConfigMap) Register(key string, item any, path ...string) *ConfigMap {
	m.locker.Lock()
	defer m.locker.Unlock()

	key = getPath(key, path...)
	_item, updated := m.tree.Insert(key, item)
	if updated {
		warnf("register key: %s updated item: %+v to %+v", key, item, _item)
	} else {
		debugf("register key: %s item: %+v", key, item)
	}
	return m
}

func (m *ConfigMap) Display() string {
	data, err := json.MarshalIndent(m.tree.ToMap(), "", "  ")
	if err != nil {
		return "json.MarshalIndent() error: " + err.Error()
	}
	return string(data)
}

func (m *ConfigMap) Load(provider Provider, parser Parser, opts ...Option) error {
	m.locker.Lock()
	defer m.locker.Unlock()

	raw, err := provider.ReadBytes()
	if err != nil {
		return err
	}

	var (
		options = newOptions(opts...)
		configs = make(map[string]any)
	)

	tagName := parser.Name()

	err = parser.Unmarshal(raw, &configs)
	if err != nil {
		return err
	}
	debugf("%s parsed config: %v", tagName, configs)

	m.tree.Walk(func(key string, item interface{}) bool {
		config, exist, err := nestedMapNoCopy(configs, strings.Split(key, PathSeparator)...)
		switch {
		case !exist:
			return false
		case err != nil:
			errorf("get item key: %s error: %s", key, err)
			return false
		}

		lock, ok := item.(sync.Locker)
		if !ok {
			lock = util.NewNopLock()
		}
		lock.Lock()
		defer lock.Unlock()

		var target = item
		// UpdateEvent
		event, ok := item.(UpdateEvent)
		if ok {
			tmp, err := copystructure.Copy(item)
			if err != nil {
				errorf("copy key: %s item: %+v error: %s", key, item, err)
				return false
			}
			target = tmp
		}
		// DecodeItemFunc
		changed, err := options.DecodeItemFunc(config, target, tagName)
		switch {
		case err != nil:
			errorf("decode to key: %s item: %+v error: %s", key, item, err)
			return false
		case !changed:
			return false
		case ok:
			// user update self
			if err = event.Update(target); err != nil {
				errorf("update key: %s item: %+v error: %s", key, item, err)
				return false
			}
		}
		// NotifyEvent
		if iface, ok := item.(NotifyEvent); ok {
			iface.Changed()
		}
		return false
	})

	return nil
}

func (m *ConfigMap) GetRaw(path ...string) (any, bool) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	key := strings.Join(path, PathSeparator)
	item, found := m.tree.Get(key)
	if !found {
		return nil, false
	}

	return item, true
}

func (m *ConfigMap) MustGetRaw(path ...string) any {
	item, ok := m.GetRaw(path...)
	if !ok {
		panic(fmt.Errorf("not found config item in path: %s", strings.Join(path, PathSeparator)))
	}
	return item
}

func (m *ConfigMap) Get(path ...string) (any, bool) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	key := strings.Join(path, PathSeparator)
	item, found := m.tree.Get(key)
	if !found {
		return nil, false
	}

	out, err := copystructure.Copy(item)
	if err != nil {
		errorf("copy key: %s item: %+v error: %s", key, item, err)
		return nil, false
	}
	if ptrOut, ok := out.(*interface{}); ok {
		return *ptrOut, true
	}
	return out, true
}

func (m *ConfigMap) MustGet(path ...string) any {
	item, ok := m.Get(path...)
	if !ok {
		panic(fmt.Errorf("not found config item in path: %s", strings.Join(path, PathSeparator)))
	}
	return item
}
