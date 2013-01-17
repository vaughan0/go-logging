package logging

import (
	"log/syslog"
)

// SyslogOutputter implements Outputter by logging to the system log daemon.
type SyslogOutputter struct {
	Writer    *syslog.Writer
	Formatter Formatter
}

// Creates a new SyslogOutputter with a custom facility (see syslog.Priority).
func NewSyslogFacility(format Formatter, tag string, facility syslog.Priority) (*SyslogOutputter, error) {
	writer, err := syslog.New(facility, tag)
	if err != nil {
		return nil, err
	}
	return &SyslogOutputter{
		Writer:    writer,
		Formatter: format,
	}, nil
}

// Creates a new SyslogOutputter with a name (tag) and the USER facility.
func NewSyslog(format Formatter, tag string) (*SyslogOutputter, error) {
	return NewSyslogFacility(format, tag, syslog.LOG_USER)
}

// Implements Outputter.
func (s SyslogOutputter) Output(msg *Message) {
	str := s.Formatter.Format(msg)
	switch msg.Level {
	case Fatal:
		s.Writer.Crit(str)
	case Error:
		s.Writer.Err(str)
	case Warn:
		s.Writer.Warning(str)
	case Notice:
		s.Writer.Notice(str)
	case Info:
		s.Writer.Info(str)
	case Debug, Trace:
		s.Writer.Debug(str)
	default:
		s.Writer.Notice(str)
	}
}
