package syslog

import (
	"errors"
	"github.com/vaughan0/go-logging"
	"log/syslog"
	"strings"
)

// SyslogOutputter implements Outputter by logging to the system log daemon.
type SyslogOutputter struct {
	Writer    *syslog.Writer
	Formatter logging.Formatter
}

// Creates a new SyslogOutputter with a custom facility (see syslog.Priority).
func NewSyslogFacility(format logging.Formatter, tag string, facility syslog.Priority) (*SyslogOutputter, error) {
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
func NewSyslog(format logging.Formatter, tag string) (*SyslogOutputter, error) {
	return NewSyslogFacility(format, tag, syslog.LOG_USER)
}

// Implements Outputter.
func (s SyslogOutputter) Output(msg *logging.Message) {
	str := s.Formatter.Format(msg)
	switch msg.Level {
	case logging.Fatal:
		s.Writer.Crit(str)
	case logging.Error:
		s.Writer.Err(str)
	case logging.Warn:
		s.Writer.Warning(str)
	case logging.Notice:
		s.Writer.Notice(str)
	case logging.Info:
		s.Writer.Info(str)
	case logging.Debug, logging.Trace:
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

var syslogPlugin = logging.OutputPluginFunc(func(options map[string]string) (result logging.Outputter, err error) {

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

	return NewSyslogFacility(logging.NewBasicFormatter(format), tag, facility)
})

func init() {
	logging.RegisterOutputPlugin("syslog", syslogPlugin)
}
