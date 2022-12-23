package event

var MAX_QUEUE_SIZE = 128
var MAX_WORKER_SIZE = 16

var MainManager *Manager

func init() {
	MainManager = NewManger(MAX_QUEUE_SIZE, MAX_WORKER_SIZE)
}
