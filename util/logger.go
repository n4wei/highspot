package util

type Logger interface {
	Printf(format string, v ...interface{})
	SetPrefix(prefix string)
}
