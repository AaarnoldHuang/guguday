package Module

import (
	"github.com/BurntSushi/toml"
	"path/filepath"
	"sync"
)

var (
	cfg  *tomlconfig
	once sync.Once
)

func Config() *tomlconfig {
	once.Do(func() {
		filePath, err := filepath.Abs("./Config/config.toml")
		if err != nil {
			panic(err)
		}
		if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
			panic(err)
		}
	})
	return cfg
}
