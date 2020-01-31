package plugin

import (
	"context"
	"fmt"
	goplugin "plugin"

	"github.com/pkg/errors"
)

type soPlugin struct {
	path   string
	plugin *goplugin.Plugin
}

func (s *soPlugin) Lookup(name string) (call Call, err error) {
	symbol, err := s.plugin.Lookup(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to lookup so symbol, name:%s path:%s", name, s.path)
	}
	call, ok := symbol.(func(ctx context.Context, req interface{}) (reply interface{}, err error))
	if !ok {
		return nil, fmt.Errorf("failed to assert so symbol, name:%s path:%s", name, s.path)
	}
	return call, nil
}

func (s *soPlugin) String() string {
	return fmt.Sprintf("type:soPlugin path:%s", s.path)
}

func (s *soPlugin) Close() error {
	return nil
}

func NewSoPlugin(path string) Plugin {
	plugin, err := goplugin.Open(path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to load so plugin, path:%s", path))
	}
	return &soPlugin{
		path:   path,
		plugin: plugin,
	}
}
