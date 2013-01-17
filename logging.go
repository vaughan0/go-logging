package logging

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// A Level describes the priority of a log message. Higher values are considered to have a higher priority. The default
// levels are negative, except for Fatal (the highest) which is zero.
type Level int

const (
	Fatal  = -iota * 100 // Unrecoverable error.
	Error                // Error condition, but possibly recoverable.
	Warn                 // Warning condition, program can still operate.
	Notice               // Normal but significant condition.
	Info                 // Informational message.
	Debug                // Debug-level message.
	Trace                // More verbose debug-level message.
)

var levelStrings = map[Level]string{
	Fatal:  "FATAL",
	Error:  "ERROR",
	Warn:   "WARN",
	Notice: "NOTICE",
	Info:   "INFO",
	Debug:  "DEBUG",
	Trace:  "TRACE",
}

// Returns a string representation of the Level, in uppercase.
func (l Level) String() string {
	if s := levelStrings[l]; s != "" {
		return s
	}
	return fmt.Sprintf("LEVEL:%d", l)
}

// A Message contains information about a logging event.
type Message struct {
	// The priority of the message.
	Level Level
	// The string part of the message, as passed by the user when the log statement was called.
	Msg string
	// The time the message was logged.
	Time time.Time
	// The name of the file where the logging statement originated.
	File string
	// The line number in the file where the logging statement originated.
	Line int
	// The Logger which logged the message.
	Logger *Logger
}

// An Outputter is responsible for logging a message to some destination.
type Outputter interface {
	Output(msg *Message)
}

type OutputterFunc func(msg *Message)

func (o OutputterFunc) Output(msg *Message) {
	o(msg)
}

// A Formatter is responsible for converting a Message into a string representation. See BasicFormatter.
type Formatter interface {
	Format(msg *Message) string
}

// Loggers are the point-of-entry for logging events.
type Logger struct {
	// The full name of the logger.
	Name string
	// The minimum level a log message can have to be logged.
	Threshold Level
	// If true, log messages will not be propagated to the parent Logger's outputs. If false, log messages will be sent up
	// the hierarchy until a Logger is found with the NoPropagate property set to true.
	NoPropagate bool
	parent      *Logger
	children    map[string]*Logger
	outputs     []Outputter
}

func newLogger(name string, parent *Logger, threshold Level) *Logger {
	return &Logger{
		Name:      name,
		Threshold: threshold,
		parent:    parent,
		children:  make(map[string]*Logger),
	}
}

func (l *Logger) log(level Level, msgstr string, stack int) {
	msg := &Message{
		Level:  level,
		Msg:    msgstr,
		Time:   time.Now(),
		Logger: l,
	}
	_, msg.File, msg.Line, _ = runtime.Caller(stack)
	l.doLog(msg)
}

func (l *Logger) doLog(msg *Message) {
	for _, output := range l.outputs {
		output.Output(msg)
	}
	if !l.NoPropagate && l.parent != nil {
		l.parent.doLog(msg)
	}
}

// Adds an Outputter to the Logger. Subsequent Messages that exceed the logger's Threshold will be sent to the
// Outputter.
func (l *Logger) AddOutput(o Outputter) {
	l.outputs = append(l.outputs, o)
}

/* Global logger hierarchy */

var lock sync.Mutex

// The root Logger. This is the ancestor of all loggers.
var Root = newLogger("root", nil, Info)

// Returns a Logger instance for the given logger name. A logger name consists of dot-separated parts, and is the basis
// of the logger hierarchy. When loggers are created (implicitly by Get) they inherit their Threshold from
func Get(fullname string) *Logger {
	lock.Lock()
	defer lock.Unlock()
	// Go down the hierarchy, creating loggers where needed
	parts := strings.Split(fullname, ".")
	logger := Root
	for _, part := range parts {
		child := logger.children[part]
		if child == nil {
			child = newLogger(fullname, logger, logger.Threshold)
			logger.children[part] = child
		}
		logger = child
	}
	return logger
}

/* Logging methods */

func (l *Logger) Log(level Level, msgstr string) {
	if l.Threshold > level {
		return
	}
	l.log(level, msgstr, 2)
}
func (l *Logger) Logf(level Level, format string, args ...interface{}) {
	if l.Threshold > level {
		return
	}
	l.log(level, fmt.Sprintf(format, args...), 2)
}

func (l *Logger) Fatal(msg string) {
	if l.Threshold > Fatal {
		return
	}
	l.log(Fatal, msg, 2)
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.Threshold > Fatal {
		return
	}
	l.log(Fatal, fmt.Sprintf(format, args...), 2)
}
func (l *Logger) Error(msg string) {
	if l.Threshold > Error {
		return
	}
	l.log(Error, msg, 2)
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.Threshold > Error {
		return
	}
	l.log(Error, fmt.Sprintf(format, args...), 2)
}
func (l *Logger) Warn(msg string) {
	if l.Threshold > Warn {
		return
	}
	l.log(Warn, msg, 2)
}
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.Threshold > Warn {
		return
	}
	l.log(Warn, fmt.Sprintf(format, args...), 2)
}
func (l *Logger) Notice(msg string) {
	if l.Threshold > Notice {
		return
	}
	l.log(Notice, msg, 2)
}
func (l *Logger) Noticef(format string, args ...interface{}) {
	if l.Threshold > Notice {
		return
	}
	l.log(Notice, fmt.Sprintf(format, args...), 2)
}
func (l *Logger) Info(msg string) {
	if l.Threshold > Info {
		return
	}
	l.log(Info, msg, 2)
}
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.Threshold > Info {
		return
	}
	l.log(Info, fmt.Sprintf(format, args...), 2)
}
func (l *Logger) Debug(msg string) {
	if l.Threshold > Debug {
		return
	}
	l.log(Debug, msg, 2)
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Threshold > Debug {
		return
	}
	l.log(Debug, fmt.Sprintf(format, args...), 2)
}
func (l *Logger) Trace(msg string) {
	if l.Threshold > Trace {
		return
	}
	l.log(Trace, msg, 2)
}
func (l *Logger) Tracef(format string, args ...interface{}) {
	if l.Threshold > Trace {
		return
	}
	l.log(Trace, fmt.Sprintf(format, args...), 2)
}
