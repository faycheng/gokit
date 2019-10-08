package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Loader interface {
	Load() (Plugin, error)
}

type loader struct {
	path string
}

func NewLoader(path string) Loader {
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &loader{
		path: path,
	}
}

// TODO: wrap error
func (l *loader) Load() (plugin Plugin, err error) {
	fd, err := os.Open(fmt.Sprintf("%s/plugin.json", l.path))
	if err != nil {
		return
	}
	defer fd.Close()
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return
	}
	config := new(PluginConfig)
	err = json.Unmarshal(content, config)
	if err != nil {
		return
	}
	config.Path = l.path
	plugin = NewPlugin(config)
	return
}
