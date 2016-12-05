package loggingmiddleware

import "log"

type options struct {
	logger func(RequestInfo)
}

func defaultOptions() *options {
	return &options{
		logger: defaultLogger,
	}
}

func WithLogger(logger func(RequestInfo)) func(*options) {
	return func(o *options) {
		o.logger = logger
	}
}

func defaultLogger(info RequestInfo) {
	log.Printf("%+v", info)
}
