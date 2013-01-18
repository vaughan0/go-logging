package logging

import (
	"errors"
	"log/syslog"
	"strings"
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

var facilityMap = map[string]syslog.Priority{
	"kern":     syslog.LOG_KERN,
	"user":     syslog.LOG_USER,
	"mail":     syslog.LOG_MAIL,
	"daemon":   syslog.LOG_DAEMON,
	"auth":     syslog.LOG_AUTH,
	"syslog":   syslog.LOG_SYSLOG,
	"lpr":      syslog.LOG_LPR,
	"news":     syslog.LOG_NEWS,
	"uucp":     syslog.LOG_UUCP,
	"cron":     syslog.LOG_CRON,
	"authpriv": syslog.LOG_AUTHPRIV,
	"ftp":      syslog.LOG_FTP,
	"local0":   syslog.LOG_LOCAL0,
	"local1":   syslog.LOG_LOCAL1,
	"local2":   syslog.LOG_LOCAL2,
	"local3":   syslog.LOG_LOCAL3,
	"local4":   syslog.LOG_LOCAL4,
	"local5":   syslog.LOG_LOCAL5,
	"local6":   syslog.LOG_LOCAL6,
	"local7":   syslog.LOG_LOCAL7,
}

var syslogPlugin = OutputPluginFunc(func(options map[string]string) (result Outputter, err error) {

	// Setup formatter
	format := options["format"]
	if format == "" {
		return nil, errors.New("syslog formatting string not specified")
	}

	tag := options["tag"]
	if tag == "" {
		return nil, errors.New("syslog tag not specified")
	}

	facility := syslog.LOG_USER
	if facilityName, ok := options["facility"]; ok {
		if facility, ok = facilityMap[strings.ToLower(facilityName)]; !ok {
			return nil, errors.New("invalid syslog facility: " + facilityName)
		}
	}

	return NewSyslogFacility(NewBasicFormatter(format), tag, facility)
})

func init() {
	RegisterOutputPlugin("syslog", syslogPlugin)
}
