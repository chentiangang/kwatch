package xlog

import "os"

type XConsole struct {
	*XLogBase
}

func NewXConsole(level int, module string) XLog {
	logger := &XConsole{}
	logger.XLogBase = &XLogBase{
		level:  level,
		module: module,
	}
	return logger
}

func (c *XConsole) Init() error {
	return nil
}

func (c *XConsole) logLevel(level int, format string, args ...interface{}) {
	if c.level > level {
		return
	}
	logData := c.formatLogger(level, c.module, format, args...)
	c.writeLog(os.Stdout, logData)
}

func (c *XConsole) LogDebug(format string, args ...interface{}) {
	c.logLevel(XLogLevelDebug, format, args...)
}

func (c *XConsole) LogTrace(format string, args ...interface{}) {
	c.logLevel(XLogLevelTrace, format, args...)
}

func (c *XConsole) LogInfo(format string, args ...interface{}) {
	c.logLevel(XLogLevelInfo, format, args...)
}

func (c *XConsole) LogWarn(format string, args ...interface{}) {
	c.logLevel(XLogLevelWarn, format, args...)
}

func (c *XConsole) LogError(format string, args ...interface{}) {
	c.logLevel(XLogLevelError, format, args...)
}

func (c *XConsole) LogFatal(format string, args ...interface{}) {
	c.logLevel(XLogLevelFatal, format, args...)
}

func (c *XConsole) SetLevel(level int) {
	c.level = level
}

func (c *XConsole) Close() {

}
