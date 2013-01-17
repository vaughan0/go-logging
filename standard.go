package logging

import (
	"bytes"
	"io"
	"path"
	"regexp"
	"strconv"
	"time"
)

// A StringWriter writes preformatted strings. StringWriters are intended to be used by StringOutputter.
type StringWriter interface {
	Write(str string)
}

// StringOutputter implements Outputter by combining a Formatter and a StringWriter.
type StringOutputter struct {
	Formatter Formatter
	Writer    StringWriter
}

// Output formats the message with the StringOutputter's Formatter, and then writes the result to the StringWriter.
func (s StringOutputter) Output(msg *Message) {
	s.Writer.Write(s.Formatter.Format(msg))
}

// IOWriter implements StringWriter by writing to an io.Writer.
type IOWriter struct {
	Writer io.Writer
}

// Implements StringWriter.
func (w IOWriter) Write(str string) {
	io.WriteString(w.Writer, str)
}

// ThresholdOutputter wraps an Outputter and only forwards messages that meet a certain threshold level.
type ThresholdOutputter struct {
	Threshold Level
	Outputter Outputter
}

func (t ThresholdOutputter) Output(msg *Message) {
	if t.Threshold > msg.Level {
		return
	}
	t.Outputter.Output(msg)
}

// BasicFormatter uses simple string templates to format messages.
type BasicFormatter struct {
	// Map of variable name to date format strings, as accepted by the Format method of time.Time objects. By default
	// contains the keys "date" (just the date), "time" (just the time), and "datetime" (date and time).
	DateVars map[string]string
	template []templatePart
}

var templateRegex = regexp.MustCompile(`^(\$[a-zA-Z]+|\$\$|[^\$]+)`)

// Returns a new BasicFormatter that uses the given template. The template may contain variables in the form $name,
// as well as arbitrary text. Variables will be substituted for their values in the result of Format. The default
// variables are:
//		level     The level of the message, as returned by Level.String.
//		msg       The string associated with the message.
//		file      The name of the file where the logging statement originated.
//		line      The line number where the logging statement originated.
//		logger    The name of the logger which was used to log the message.
// Variables from DateVars ($date, $time and $datetime, by default) are also included.
//
// For example: If the template is "[$level] $time - $msg\n", then the call logger.Warn("oh no!") could produce
// an output of: "[WARN] 15:04:05 - oh no!\n".
func NewBasicFormatter(template string) *BasicFormatter {
	parts := []templatePart{}
	remain := template
	for len(remain) > 0 {
		match := templateRegex.FindString(remain)
		if match == "" {
			panic("invalid template: " + template)
		}
		switch {
		case match == "$$":
			parts = append(parts, templatePart{"$", false})
		case match[0] == '$':
			parts = append(parts, templatePart{match[1:], true})
		default:
			parts = append(parts, templatePart{match, false})
		}
		remain = remain[len(match):]
	}
	return &BasicFormatter{
		template: parts,
		DateVars: map[string]string{
			"date":     "02/01/2006",
			"time":     "15:04:05",
			"datetime": time.ANSIC,
		},
	}
}

// Implements Formatter.
func (b *BasicFormatter) Format(msg *Message) string {
	vars := b.getVars(msg)
	var result bytes.Buffer
	for _, part := range b.template {
		if part.Var {
			result.WriteString(vars[part.Str])
		} else {
			result.WriteString(part.Str)
		}
	}
	return result.String()
}

func (b *BasicFormatter) getVars(msg *Message) map[string]string {
	vars := map[string]string{
		"level":  msg.Level.String(),
		"msg":    msg.Msg,
		"file":   path.Base(msg.File),
		"line":   strconv.Itoa(msg.Line),
		"logger": msg.Logger.Name,
	}
	for key, layout := range b.DateVars {
		vars[key] = msg.Time.Format(layout)
	}
	return vars
}

type templatePart struct {
	Str string
	Var bool
}
