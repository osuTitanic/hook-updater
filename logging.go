package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	QUIET   int = 60
	ANOMALY int = 50
	ERROR   int = 40
	WARNING int = 30
	INFO    int = 20
	DEBUG   int = 10
	VERBOSE int = 0
)

const (
	ColorReset  = "\033[0m"
	ColorBlack  = "\033[30m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGrey   = "\033[0;90m"
)

const (
	ColorBoldBlack  = "\033[1;30m"
	ColorBoldRed    = "\033[1;31m"
	ColorBoldGreen  = "\033[1;32m"
	ColorBoldYellow = "\033[1;33m"
	ColorBoldBlue   = "\033[1;34m"
	ColorBoldPurple = "\033[1;35m"
	ColorBoldCyan   = "\033[1;36m"
	ColorBoldWhite  = "\033[1;37m"
)

type Logger struct {
	logger *log.Logger
	name   string
	level  int
}

func CreateLogger(name string, level int) *Logger {
	l := log.New(os.Stdout, "", 0)
	return &Logger{
		logger: l,
		name:   name,
		level:  level,
	}
}

func (c *Logger) SetLevel(level int) {
	c.level = level
}

func (c *Logger) GetLevel() int {
	return c.level
}

func (c *Logger) SetName(name string) {
	c.name = name
}

func (c *Logger) GetName() string {
	return c.name
}

func (c *Logger) Info(msg ...any) {
	if c.level > INFO {
		return
	}
	c.logger.Println(c.formatLogMessage("INFO", concatMessage(msg...)))
}

func (c *Logger) Infof(format string, msg ...any) {
	if c.level > INFO {
		return
	}
	c.logger.Println(c.formatLogMessage("INFO", fmt.Sprintf(format, msg...)))
}

func (c *Logger) Error(msg ...any) {
	if c.level > ERROR {
		return
	}
	c.logger.Println(c.formatLogMessage("ERROR", concatMessage(msg...)))
}

func (c *Logger) Errorf(format string, msg ...any) {
	if c.level > ERROR {
		return
	}
	c.logger.Println(c.formatLogMessage("ERROR", fmt.Sprintf(format, msg...)))
}

func (c *Logger) Warning(msg ...any) {
	if c.level > WARNING {
		return
	}
	c.logger.Println(c.formatLogMessage("WARNING", concatMessage(msg...)))
}

func (c *Logger) Warningf(format string, msg ...any) {
	if c.level > WARNING {
		return
	}
	c.logger.Println(c.formatLogMessage("WARNING", fmt.Sprintf(format, msg...)))
}

func (c *Logger) Debug(msg ...any) {
	if c.level > DEBUG {
		return
	}
	c.logger.Println(c.formatLogMessage("DEBUG", concatMessage(msg...)))
}

func (c *Logger) Debugf(format string, msg ...any) {
	if c.level > DEBUG {
		return
	}
	c.logger.Println(c.formatLogMessage("DEBUG", fmt.Sprintf(format, msg...)))
}

func (c *Logger) Verbose(msg ...any) {
	if c.level > VERBOSE {
		return
	}
	c.logger.Println(c.formatLogMessage("VERBOSE", concatMessage(msg...)))
}

func (c *Logger) Verbosef(format string, msg ...any) {
	if c.level > VERBOSE {
		return
	}
	c.logger.Println(c.formatLogMessage("VERBOSE", fmt.Sprintf(format, msg...)))
}

func (c *Logger) Anomaly(msg ...any) {
	if c.level > ANOMALY {
		return
	}
	c.logger.Println(c.formatLogMessage("ANOMALY", concatMessage(msg...)))
}

func (c *Logger) Anomalyf(format string, msg ...any) {
	if c.level > ANOMALY {
		return
	}
	c.logger.Println(c.formatLogMessage("ANOMALY", fmt.Sprintf(format, msg...)))
}

func (c *Logger) formatLogMessage(level string, msg string) string {
	color := getColor(level)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(
		"[%s] - <%s> %s%s: %s%s",
		timestamp,
		c.name,
		color,
		level,
		msg,
		ColorReset,
	)
}

func getColor(level string) string {
	switch level {
	case "INFO":
		return ColorBlue
	case "ERROR":
		return ColorRed
	case "ANOMALY":
		return ColorBoldYellow
	case "WARNING":
		return ColorYellow
	case "DEBUG":
		return ColorWhite
	case "VERBOSE":
		return ColorGrey
	default:
		return ColorReset
	}
}

func concatMessage(msg ...any) string {
	// Concatenate all strings in msg
	log := ""

	for _, m := range msg {
		log += fmt.Sprintf("%v ", m)
	}

	return log
}
