package eventbus

import (
	"errors"
	"fmt"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"
)

type handlerRec struct {
	name    string
	fn      HandlerFunc
	once    bool // if true, auto-unregister after first successful run
	channel string
}

type task struct {
	ev *Event
	h  handlerRec
}

// bus implements Bus.
type bus struct {
	// configuration
	workerCount int
	queueSize   int

	// workers
	queues []chan task
	wg     sync.WaitGroup

	// lifecycle
	started   atomic.Bool
	stopping  atomic.Bool
	stopOnce  sync.Once
	stopCh    chan struct{}
	drainedCh chan struct{}

	// routing & bookkeeping
	mu       sync.RWMutex
	handlers map[string]map[string]handlerRec // eventId -> handlerName -> handlerRec
	pending  []*Event                         // events published before Start()

	// rendezvous for Call/EchoId
	waitMu     sync.Mutex
	echoWaiter map[string]chan *Event
	// simple id source for EchoId if caller doesn't provide
	idCtr atomic.Uint64

	// logger
	log Logger
}

// New creates a new Bus.
// workerCount >= 1, queueSize >= 1.
func New(opts ...Option) Bus {
	option := options{
		log:        Log,
		workerSize: 10,
		queueSize:  100,
	}
	for _, opt := range opts {
		opt(&option)
	}
	b := &bus{
		workerCount: option.workerSize,
		queueSize:   option.queueSize,
		queues:      make([]chan task, option.workerSize),
		stopCh:      make(chan struct{}),
		drainedCh:   make(chan struct{}),
		handlers:    make(map[string]map[string]handlerRec),
		pending:     make([]*Event, 0, 16),
		echoWaiter:  make(map[string]chan *Event),
		log:         option.log,
	}
	for i := 0; i < option.workerSize; i++ {
		q := make(chan task, option.queueSize)
		b.queues[i] = q
		go b.workerLoop(q)
	}
	return b
}

func (b *bus) workerLoop(q chan task) {
	for {
		select {
		case <-b.stopCh:
			// Drain quickly without executing tasks (immediate stop).
			for {
				select {
				case <-q:
					// drop
				default:
					return
				}
			}
		case t := <-q:
			func() {
				defer func() {
					if r := recover(); r != nil {
						b.log.Printf("handler panic recovered: event=%s handler=%s panic=%v", t.ev.Id, t.h.name, r)
					}
					b.wg.Done()
				}()
				// Execute handler
				t.h.fn(t.ev)
				// If it was a once-handler, unregister it after execution.
				if t.h.once {
					_ = b.Unsubscribe(t.ev.Id, t.h.name)
				}
			}()
		}
	}
}

func (b *bus) Start() error {
	if b.started.Swap(true) {
		return nil
	}
	// Flush pending
	b.mu.Lock()
	pending := b.pending
	b.pending = nil
	b.mu.Unlock()

	for _, ev := range pending {
		err := b.PublishEvent(ev)
		if err != nil {
			b.log.Printf("failed to publish event: %v", err)
		}
	}
	return nil
}

func (b *bus) Wait() error {
	// Wait for all in-flight tasks (that were queued before this call) to finish.
	// If Stop() has been called, Wait returns after workers exit.
	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-b.drainedCh:
		// Stopped
		<-done
		return nil
	}
}

func (b *bus) Stop() error {
	b.stopOnce.Do(func() {
		b.stopping.Store(true)
		close(b.stopCh)    // signal workers to stop immediately
		close(b.drainedCh) // allow Wait() to proceed
	})
	return nil
}

func (b *bus) Subscribe(channel string, eventId, handlerName string, fn HandlerFunc) error {
	if eventId == "" || handlerName == "" || fn == nil {
		return errors.New("invalid Subscribe args")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	m := b.handlers[eventId]
	if m == nil {
		m = make(map[string]handlerRec)
		b.handlers[eventId] = m
	}
	m[handlerName] = handlerRec{name: handlerName, fn: fn, channel: channel}
	return nil
}

func (b *bus) SubscribeAny(eventId, handlerName string, fn HandlerFunc) error {
	if eventId == "" || handlerName == "" || fn == nil {
		return errors.New("invalid Subscribe args")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	m := b.handlers[eventId]
	if m == nil {
		m = make(map[string]handlerRec)
		b.handlers[eventId] = m
	}
	m[handlerName] = handlerRec{name: handlerName, fn: fn, channel: ""}
	return nil
}

func (b *bus) SubscribeOnce(eventId, handlerName string, fn HandlerFunc) error {
	if eventId == "" || handlerName == "" || fn == nil {
		return errors.New("invalid SubscribeOnce args")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	m := b.handlers[eventId]
	if m == nil {
		m = make(map[string]handlerRec)
		b.handlers[eventId] = m
	}
	m[handlerName] = handlerRec{name: handlerName, fn: fn, once: true}
	return nil
}

func (b *bus) Unsubscribe(eventId, handlerName string) error {
	if eventId == "" || handlerName == "" {
		return errors.New("invalid Unsubscribe args")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if m := b.handlers[eventId]; m != nil {
		delete(m, handlerName)
		if len(m) == 0 {
			delete(b.handlers, eventId)
		}
	}
	return nil
}

func (b *bus) Publish(eventId string, data interface{}) error {
	return b.PublishEvent(&Event{Id: eventId, Data: data})
}

func (b *bus) PublishEvent(ev *Event) error {
	if ev == nil || ev.Id == "" {
		return errors.New("invalid PublishEvent args")
	}
	// If stopping, drop events.
	if b.stopping.Load() {
		return errors.New("bus is stopping")
	}

	// Rendezvous: if this looks like a reply (EchoId set) and someone is waiting, deliver.
	if ev.EchoId != "" {
		b.waitMu.Lock()
		if ch, ok := b.echoWaiter[ev.Id+ev.EchoId]; ok {
			select {
			case ch <- ev:
			default:
			}
		}
		b.waitMu.Unlock()
	}

	b.mu.RLock()
	started := b.started.Load()
	if !started {
		// queue as pending (publish-before-start)
		b.mu.RUnlock()
		b.mu.Lock()
		if !b.started.Load() {
			b.pending = append(b.pending, cloneEvent(ev))
			b.mu.Unlock()
			return nil
		}
		// started flipped while acquiring lock, fallthrough to publish now
		b.mu.Unlock()
		b.mu.RLock()
	}
	// Snapshot handlers for this event id.
	m := b.handlers[ev.Id]
	if len(m) == 0 {
		b.mu.RUnlock()
		return nil
	}

	// Make a stable copy to avoid holding the lock during execution.
	hs := make([]handlerRec, 0, len(m))
	for _, h := range m {
		if ev.Channel != "" && h.channel != "" && ev.Channel != h.channel {
			// channel not match pass
			continue
		}
		hs = append(hs, h)
	}
	b.mu.RUnlock()

	// Enqueue each handler on its shard (worker) based on (eventId, handlerName).
	for _, h := range hs {
		idx := shardIndex(b.workerCount, ev.Id, h.name)
		b.wg.Add(1)
		select {
		case b.queues[idx] <- task{ev: cloneEvent(ev), h: h}:
		default:
			// Backpressure: if shard queue is full, block (ensures ordering) but still bounded overall.
			b.queues[idx] <- task{ev: cloneEvent(ev), h: h}
		}
	}
	return nil
}

// Call publishes a request and waits for a response event with the same EchoId.
// NOTE: Handlers should reply by publishing an Event with the SAME EchoId.
// Use Reply helper below.
func (b *bus) Call(eventId string, data interface{}, subEvtId string) (*Event, error) {
	if eventId == "" {
		return nil, errors.New("empty eventId")
	}
	echo := b.nextEchoId()
	wait := make(chan *Event, 1)

	b.waitMu.Lock()
	b.echoWaiter[subEvtId+echo] = wait
	b.waitMu.Unlock()
	defer func() {
		b.waitMu.Lock()
		delete(b.echoWaiter, subEvtId+echo)
		b.waitMu.Unlock()
	}()

	b.PublishEvent(&Event{Id: eventId, EchoId: echo, Data: data})

	timeout := time.After(6 * time.Second)

	// No timeout specified in interface; block until reply or stop.
	select {
	case resp := <-wait:
		return resp, nil
	case <-timeout:
		return nil, errors.New("call timeout")
	case <-b.drainedCh:
		return nil, errors.New("bus stopped")
	}
}

func (b *bus) nextEchoId() string {
	x := b.idCtr.Add(1)
	return fmt.Sprintf("echo-%d-%d", time.Now().UnixNano(), x)
}

func shardIndex(n int, eventId, handlerName string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(eventId))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(handlerName))
	return int(h.Sum32() % uint32(n))
}

func cloneEvent(e *Event) *Event {
	if e == nil {
		return nil
	}
	// shallow clone is fine; Data is user-owned
	cp := *e
	return &cp
}
