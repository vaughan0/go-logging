package logging

import (
	"os"
	"io"
	"errors"
	"strconv"
)

// A WriterPlugin implements OutputPlugin by using a function to choose an io.Writer.
type WriterPlugin func(options map[string]string) (writer io.Writer, err error)

func (chooser WriterPlugin) CreateOutputter(options map[string]string) (result Outputter, err error) {

	// Setup formatter
	format := options["format"]
	if format == "" {
		return nil, errors.New("console formatting string not specified")
	}
	format += "\n"
	formatter := NewBasicFormatter(format)

	// Determine output stream to use
	output, err := chooser(options)
	if err != nil {
		return
	}

	return StringOutputter{
		Writer: IOWriter{output},
		Formatter: formatter,
	}, nil
}

var consolePlugin = WriterPlugin(func(options map[string]string) (output io.Writer, err error) {
	stream := options["stream"]
	switch {
	case stream == "stdout":
		output = os.Stdout
	case stream == "stderr":
		output = os.Stderr
	case stream == "":
		err = errors.New("console stream not specified")
	default:
		if fd, err := strconv.Atoi(stream); err == nil {
			output = os.NewFile(uintptr(fd), "logging_output")
		} else {
			err = errors.New("invalid console stream: " + stream)
		}
	}
	return
})

var filePlugin = WriterPlugin(func(options map[string]string) (output io.Writer, err error) {
	path := options["file"]
	if path == "" {
		err = errors.New("file option not specified")
	} else {
		output, err = os.OpenFile(path, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	}
	return
})

func init() {
	RegisterOutputPlugin("console", consolePlugin)
	RegisterOutputPlugin("file", filePlugin)
}
