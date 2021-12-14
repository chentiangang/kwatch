package xlog

const (
	XLogLevelDebug = iota
	XLogLevelTrace
	XLogLevelInfo
	XLogLevelWarn
	XLogLevelError
	XLogLevelFatal
)

const (
	red = uint8(iota + 91)
	green
	yellow
	blue
	magenta
)

func getLevelStr(level int) (lev string, color uint8) {
	switch level {
	case XLogLevelDebug:
		return "DEBUG", blue
	case XLogLevelTrace:
		return "TRACE", blue
	case XLogLevelInfo:
		return "INFO", green
	case XLogLevelWarn:
		return "WARN", yellow
	case XLogLevelError:
		return "ERROR", red
	case XLogLevelFatal:
		return "FATAL", red
	default:
		return "UNKNOWN", magenta
	}
}
