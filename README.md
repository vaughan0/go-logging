go-logging
==========

Simple logging library for Go (golang).

View the API documentation [here](http://godoc.org/github.com/vaughan0/go-logging).

Getting Started
---------------

Configure the logging library:

```go
import "github.com/vaughan0/go-logging"

logging.DefaultSetup()
```

Get a logger and log some messages:

```go
log := logging.Get("my.logger")

log.Info("This is an informational message")
log.Warn("This is a WARNING")
log.Warnf("This message was made from a %s", "format string")
```

This should produce output similar to the following:

```
[INFO] Fri Jan 18 19:18:01 2013 - This is an informational message
[WARN] Fri Jan 18 19:18:01 2013 - This is a WARNING
[WARN] Fri Jan 18 19:18:01 2013 - This message was made from a format string
```

go-logging defines several logging "levels", which represent the priority of a message. The default levels are, in
ascending order of priority:
* Trace
* Debug
* Info
* Notice
* Warn
* Error
* Fatal

Configuring
-----------

go-logging is configured by using a simple INI file. The easiest way to do this is to call logging.MustSetup(), which
will parse the INI file specified by the GO_LOGGING_CONFIG environment variable.

Here is an example configuration file:

```ini
# Configure logger thresholds and outputs in the "loggers" section:
[loggers]
# By default, only messages with a level of INFO or higher will be logged.
# The default output is the "console" output, which will be defined in the next section.
root = INFO, console
# Turn off all messages except FATAL ones from any loggers from vaughan0's libraries.
vaughan0 = FATAL

# All other sections are output definitions:
[console]
type = console
stream = stderr
format = $time $level: ($file:$line) $msg

[logfile]
type = file
file = logging-is-fun.txt
format = $time $level ($logger) $msg
```

If you load that file and run the code from "Getting Started", you will get output similar to the following:

```
19:26:46 INFO: (myfile.go:12) This is an informational message
19:26:46 WARN: (myfile.go:13) This is a WARNING
19:26:46 WARN: (myfile.go:14) This message was made from a format string
```

Logger Hierarchy
----------------

Loggers form a hierarchy and inherit their thresholds from their parent loggers, unless they have been overridden by the
configuration. The hierarchy is formed by splitting logger names up by full stops, ie. the name "foo.bar.baz" refers to
the "baz" logger, whose parent is "bar", whose parent is "foo".

Logger outputs are also inherited, however if outputs are defined for say, the "A.B" logger, messages will still be
sent to A's outputs _as well as_ B's outputs. This behaviour can be undesirable and may be disabled on a per-logger
basis by using the "nopropagate" option.

Example of nopropagate:

```ini
[loggers]
# A's messages (with levels of INFO or higher) will be sent to the "console" output
A = INFO, console
# B's messages will be sent exclusively to the "special" output
A.B = TRACE, special, nopropagate
# All of A.C's messages (of any level) will be sent to the "console" output
A.C = TRACE
```
