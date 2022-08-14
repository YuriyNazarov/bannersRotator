package app

type Logger interface {
	Close()
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}
