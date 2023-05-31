package main

import (
	"log"
	"reflect"
	"strings"

	"github.com/ElfLabs/configmap"
	"github.com/ElfLabs/configmap/example/logger"
	"github.com/ElfLabs/configmap/parsers/json"
	"github.com/ElfLabs/configmap/parsers/yaml"
	"github.com/ElfLabs/configmap/providers/file"
)

func ConfigMapDemo() {
	conf := configmap.New()

	conf.Register("debug", true)
	conf.Register("logger", logger.NewDefaultConfig())
	conf.Register("logger.server", logger.NewDefaultConfig())
	conf.Register("logger.server", logger.NewDefaultConfig())
	conf.Register("logger.database", logger.NewDefaultConfig())
	conf.Register("database.redis", map[string]any{
		"addr": ":9673",
		"db":   1,
	})

	log.Printf("configmap: %s", conf.Display())

	log.Printf("%s", strings.Repeat("-", 120))
	if err := conf.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("load file error: %s", err)
	}
	log.Printf("config.yaml to configmap: %s", conf.Display())

	log.Printf("%s", strings.Repeat("-", 120))
	if err := conf.Load(file.Provider("config.json"), json.Parser()); err != nil {
		log.Fatalf("to load file error: %s", err)
	}
	log.Printf("config.json configmap: %s", conf.Display())

	get := func(path ...string) {
		item, _ := conf.Get(path...)
		log.Println("get key:", strings.Join(path, "."), "item:", item, "type:", reflect.TypeOf(item))
	}
	get("debug")
	get("logger")
	get("logger.server")
	get("logger.server.1")
	get("logger.database")
	get("logger.redis")
	get("database.redis")
	get("server")
	get("logger.ser")
}
