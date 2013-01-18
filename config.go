package logging

import (
	"errors"
)

var ErrTypeNotSpecified = errors.New("plugin name not specified")

type ErrUnknownPlugin string

func (e ErrUnknownPlugin) Error() string {
	return "unknown logging plugin: " + string(e)
}

type OutputPlugin interface {
	CreateOutputter(options map[string]string) (Outputter, error)
}

var outputPlugins map[string]OutputPlugin

func RegisterOutputPlugin(name string, plugin OutputPlugin) {
	lock.Lock()
	defer lock.Unlock()
	outputPlugins[name] = plugin
}

func newOutputterConfig(config map[string]string) (Outputter, error) {
	name, ok := config["type"]
	if !ok {
		return nil, ErrTypeNotSpecified
	}
	lock.Lock()
	defer lock.Unlock()
	plugin := outputPlugins[name]
	if plugin == nil {
		return nil, ErrUnknownPlugin(name)
	}
	return plugin.CreateOutputter(config)
}
