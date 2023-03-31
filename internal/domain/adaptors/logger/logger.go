package logger

// Logger interface for logging within the application
type Logger interface {
	Fatal(message string)
	Error(message string)
	Warn(message string)
	Debug(message string)
	Info(message string)
	Trace(message string)
}
