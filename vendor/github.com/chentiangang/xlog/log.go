package xlog

var logger XLog = newXLog("console", XLogLevelDebug, "", "default")

type XLog interface {
	Init() error
	LogDebug(fmt string, args ...interface{})
	LogTrace(fmt string, args ...interface{})
	LogInfo(fmt string, args ...interface{})
	LogWarn(fmt string, args ...interface{})
	LogError(fmt string, args ...interface{})
	LogFatal(fmt string, args ...interface{})

	Close()
	SetLevel(level int)
}

func newXLog(logType string, level int, filename, module string) XLog {
	var logger XLog
	switch logType {
	case "file", "File":
		logger = NewXFile(level, filename, module)
	case "console", "Console":
		logger = NewXConsole(level, module)
	default:
		logger = NewXFile(level, filename, module)
	}
	return logger
}

func Init(logType string, level string, filename, module string) error {
	var l int
	switch level {
	case "debug":
		l = XLogLevelDebug
	case "trace":
		l = XLogLevelTrace
	case "info":
		l = XLogLevelInfo
	case "warn":
		l = XLogLevelWarn
	case "error":
		l = XLogLevelError
	case "fatal":
		l = XLogLevelFatal
	default:
		l = XLogLevelDebug
	}

	logger = newXLog(logType, l, filename, module)
	return logger.Init()
}

func Debug(format string, args ...interface{}) {
	logger.LogDebug(format, args...)
}

func Trace(format string, args ...interface{}) {
	logger.LogTrace(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.LogInfo(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.LogWarn(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.LogError(format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.LogFatal(format, args...)
}

func Close() {
	logger.Close()
}

func SetLevel(level int) {
	logger.SetLevel(level)
}
