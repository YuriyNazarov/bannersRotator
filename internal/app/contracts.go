package app

type Logger interface {
	Close()
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}
