package eventbus

type options struct {
	log        Logger
	workerSize int
	queueSize  int
}

type Option func(*options)

func WithLogger(logger Logger) Option {
	return func(o *options) { o.log = logger }
}

func WithWorkerSize(workerSize int) Option {
	return func(o *options) {
		if workerSize >= 1 {
			o.workerSize = workerSize
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
