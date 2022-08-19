package app

import "errors"

type Logger interface {
	Close()
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

var ErrNoBanners = errors.New("no banners available for selected slot")
