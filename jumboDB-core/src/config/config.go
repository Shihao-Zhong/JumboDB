package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"sync"
)

type TomlConfig struct {
	Title      string
	Connection connection `toml:"connection"`
	Storage storage `toml:"storage"`
}

type connection struct {
	Port int
}

type storage struct {
	BloomFilterFalsePositiveRate float64 `toml:"bloomFilterFalsePositiveRate"`
	MemoryTableSize int `toml:"memoryTableSize"`
	SkipListLevel int `toml:"skipListLevel"`
	SSTableIndexLocation string `toml:"SSTableIndexLocation"`
	WALLocation          string `toml:"WALLocation"`
	LSMIndexLocation     string `toml:"LSMIndexLocation"`
	Engine               string `toml:"engine"`
}

var conf *TomlConfig
var Path string
var mu sync.Mutex

func GetConfig() *TomlConfig {
	mu.Lock()
	defer mu.Unlock()
	if conf == nil {
		if _, err := toml.DecodeFile(Path, &conf); err != nil {
			log.Fatalf("Unable to parse JumboDB config in path [%s] with error [%s]", Path, err)
		}
	}
	return conf
}