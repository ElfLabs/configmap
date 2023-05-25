package configmap

import (
	"encoding/json"
	"log"
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

	_item, updated := m.tree.Insert(getPath(key, path...), item)
	log.Printf("register item: %+v updated: %t", _item, updated)
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

	err = parser.Unmarshal(raw, &configs)
	if err != nil {
		return err
	}

	tagName := parser.Name()

	m.tree.Walk(func(key string, item interface{}) bool {
		config, exist, err := nestedMapNoCopy(configs, strings.Split(key, PathSeparator)...)
		switch {
		case !exist:
			return false
		case err != nil:
			log.Printf("get %s error: %s", key, err)
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
				log.Printf("copy item: %s error: %s", key, err)
				return false
			}
			target = tmp
		}

		changed, err := options.DecodeItemFunc(config, target, tagName)
		switch {
		case err != nil:
			log.Printf("decode to item: %s error: %s", key, err)
			return false
		case !changed:
			return false
		case ok:
			if err = event.Update(target); err != nil {
				log.Printf("update item: %s error: %s", key, err)
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

func (m *ConfigMap) Get(path ...string) (any, bool) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	item, found := m.tree.Get(strings.Join(path, PathSeparator))
	if !found {
		return nil, false
	}

	out, _ := copystructure.Copy(item)
	if ptrOut, ok := out.(*interface{}); ok {
		return *ptrOut, true
	}
	return out, true
}
