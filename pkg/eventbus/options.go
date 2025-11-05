package eventbus

type options struct {
	log           Logger
	maxWorkerSize int
	queueSize     int
}

type Option func(*options)

func WithLogger(logger Logger) Option {
	return func(o *options) { o.log = logger }
}

func WithMaxWorkerSize(maxWorkerSize int) Option {
	return func(o *options) {
		if maxWorkerSize >= 1 {
			o.maxWorkerSize = maxWorkerSize
		}
	}
}

func WithQueueSize(queueSize int) Option {
	return func(o *options) {
		if queueSize >= 1 {
			o.queueSize = queueSize
		}
	}
}
