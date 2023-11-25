package queue

type T interface{}

type ConcurrentQueue struct {
	Queue       chan T
	QueueSize   int
	Concurrency int
}

var (
	concurrency = 0
)

func New(queueSize int, concurrency int) *ConcurrentQueue {
	q := ConcurrentQueue{
		QueueSize:   queueSize,
		Queue:       make(chan T, queueSize),
		Concurrency: concurrency,
	}

	return &q
}

func (cq *ConcurrentQueue) Enqueue(item T) {
	cq.Queue <- item
}

func (cq *ConcurrentQueue) DispatchProcess(process func(interface{})) {
	for {
		select {
		case item := <-cq.Queue:
			if concurrency < cq.Concurrency {
				go func() {
					defer func() {
						concurrency--
					}()
					concurrency++
					process(item)
				}()
			} else {
				process(item)
			}
		}
	}
}
