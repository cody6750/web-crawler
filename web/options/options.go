package options

import (
	"time"
)

var (
	defaultLogLevel     = "INFO"
	defaultPort         = ":9090"
	defaultIdleTimeout  = time.Second * 120
	defaultReadTimeout  = time.Second * 60
	defaultWriteTimeout = time.Second * 60
)

// Options includes the overrideable option configurations for the web crawler server
type Options struct {
	LogLevel     string
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func New() *Options {
	return &Options{
		LogLevel:     defaultLogLevel,
		Port:         defaultPort,
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
}
