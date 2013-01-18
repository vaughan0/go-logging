package logging

import (
	"errors"
	"github.com/vaughan0/go-ini"
	"io"
	"os"
	"strings"
)

type Config ini.File

// ErrTypeNotSpecified is returned when an output section does not contain a "type" option.
var ErrTypeNotSpecified = errors.New("plugin name not specified")

// ErrUnknownPlugin is returned when an output section's "type" option refers to an unknown plugin.
type ErrUnknownPlugin string

func (e ErrUnknownPlugin) Error() string {
	return "unknown logging plugin: " + string(e)
}

// An OutputPlugin is responsible for creating Outputters from simple key-value configuration variables.
type OutputPlugin interface {
	CreateOutputter(options map[string]string) (Outputter, error)
}

// OutputPluginFunc is a utility type that implements OutputPlugin.
type OutputPluginFunc func(options map[string]string) (Outputter, error)

// Implements OutputPlugin.
func (o OutputPluginFunc) CreateOutputter(options map[string]string) (Outputter, error) {
	return o(options)
}

var outputPlugins = make(map[string]OutputPlugin)

// Registers an output plugin by name.
func RegisterOutputPlugin(name string, plugin OutputPlugin) {
	lock.Lock()
	defer lock.Unlock()
	outputPlugins[name] = plugin
}

// Loads the appropriate plugin and creates an outputter, given a configuration section.
func newOutputterConfig(config ini.Section) (Outputter, error) {
	// Get plugin from the "type" option
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

	output, err := plugin.CreateOutputter(config)
	if err != nil {
		return nil, err
	}

	// Check for the "threshold" option
	if thresh, ok := config["threshold"]; ok {
		if level, ok := reverseLevelStrings[strings.ToUpper(thresh)]; ok {
			output = ThresholdOutputter{level, output}
		} else {
			return nil, errors.New("invalid threshold: " + thresh)
		}
	}
	return output, nil
}

// Applies configuration to the logging hierarchy.
func (c Config) apply() (err error) {

	// Create outputters
	outputters := make(map[string]Outputter)
	for key, section := range c {
		if key != "loggers" && key != "" {
			var output Outputter
			if output, err = newOutputterConfig(section); err != nil {
				return
			}
			outputters[key] = output
		}
	}

	loggerSection := c["loggers"]
	if loggerSection == nil {
		return errors.New("loggers section not specified")
	}

	// Setup loggers
	for name, config := range loggerSection {
		parts := strings.Split(config, ",")
		level, ok := reverseLevelStrings[strings.ToUpper(parts[0])]
		if !ok {
			return errors.New("unknown logging level: " + parts[0])
		}
		// Get the logger by its name, treating "root" as a special name
		var logger *Logger
		if name == "root" {
			logger = Root
		} else {
			logger = Get(name)
		}
		logger.Threshold = level
		// Handle extra options
		for _, outputKey := range parts[1:] {
			if outputKey == "nopropagate" {
				logger.NoPropagate = true
			} else {
				// Assign an outputter
				if outputter := outputters[outputKey]; outputter != nil {
					logger.AddOutput(outputter)
				} else {
					return errors.New("unknown logging output: " + outputKey)
				}
			}
		}
	}

	Root.configure()
	configured = true
	return nil
}

// Configures the logging hierarchy from an io.Reader, which should return valid INI source code.
func SetupReader(config io.Reader) (err error) {
	file, err := ini.Load(config)
	if err != nil {
		return
	}
	return Config(file).apply()
}

// Configures the logging hierarchy from an INI file.
func SetupFile(filename string) (err error) {
	file, err := ini.LoadFile(filename)
	if err != nil {
		return
	}
	return Config(file).apply()
}

// Automatically configures the logging hierarchy by loading the INI file specified by the GO_LOGGING_CONFIG environment
// variable.
func Setup() (err error) {
	path := os.Getenv("GO_LOGGING_CONFIG")
	if path == "" {
		return errors.New("GO_LOGGING_CONFIG not set")
	}
	return SetupFile(path)
}
