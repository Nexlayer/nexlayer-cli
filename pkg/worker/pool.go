package worker

import (
	"container/heap"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Job func(ctx context.Context) error

type PoolConfig struct {
	MinWorkers    int           // Minimum number of workers
	MaxWorkers    int           // Maximum number of workers
	QueueSize     int           // Size of the job queue
	ScaleInterval time.Duration // How often to check for scaling
	IdleTimeout   time.Duration // How long a worker can be idle before shutdown
}

type Pool struct {
	config     PoolConfig
	jobs       chan Job
	results    chan error
	priorities chan priorityJob
	workers    map[int]*workerInfo
	workersWg  sync.WaitGroup
	metrics    *poolMetrics
	mu         sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type workerInfo struct {
	id        int
	lastUsed  time.Time
	processed int64
	errors    int64
	done      chan struct{}
}

type priorityJob struct {
	job      Job
	priority int
	done     chan struct{}
}

type poolMetrics struct {
	activeWorkers  int64
	completedJobs  int64
	failedJobs     int64
	queuedJobs     int64
	avgProcessTime time.Duration
}

type priorityQueue []*priorityJob

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*priorityJob)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// NewPool creates a new worker pool with auto-scaling capabilities
func NewPool(config PoolConfig) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		config:     config,
		jobs:       make(chan Job, config.QueueSize),
		priorities: make(chan priorityJob, config.QueueSize),
		results:    make(chan error, config.QueueSize),
		workers:    make(map[int]*workerInfo),
		metrics:    &poolMetrics{},
		ctx:        ctx,
		cancelFunc: cancel,
	}

	return p
}

// Start initializes the worker pool and starts the scaling manager
func (p *Pool) Start() {
	// Start minimum number of workers
	for i := 0; i < p.config.MinWorkers; i++ {
		p.startWorker()
	}

	// Start the scaling manager
	go p.manageScale()

	// Start priority job handler
	go p.handlePriorityJobs()
}

// Submit adds a job to the pool with normal priority
func (p *Pool) Submit(job Job) {
	atomic.AddInt64(&p.metrics.queuedJobs, 1)
	p.jobs <- job
}

// SubmitPriority adds a job with specified priority (higher number = higher priority)
func (p *Pool) SubmitPriority(job Job, priority int) {
	done := make(chan struct{})
	p.priorities <- priorityJob{job: job, priority: priority, done: done}
	atomic.AddInt64(&p.metrics.queuedJobs, 1)
	<-done // Wait for priority job to be queued
}

func (p *Pool) handlePriorityJobs() {
	priorityQueue := &priorityQueue{}
	heap.Init(priorityQueue)

	for {
		select {
		case pj := <-p.priorities:
			heap.Push(priorityQueue, pj)
			close(pj.done)
		case p.jobs <- heap.Pop(priorityQueue).(*priorityJob).job:
			// Job was sent to a worker
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Pool) startWorker() {
	p.mu.Lock()
	id := len(p.workers)
	worker := &workerInfo{
		id:       id,
		lastUsed: time.Now(),
		done:     make(chan struct{}),
	}
	p.workers[id] = worker
	p.mu.Unlock()

	atomic.AddInt64(&p.metrics.activeWorkers, 1)
	p.workersWg.Add(1)

	go func() {
		defer p.workersWg.Done()
		defer atomic.AddInt64(&p.metrics.activeWorkers, -1)

		for {
			select {
			case job := <-p.jobs:
				start := time.Now()
				err := job(p.ctx)

				worker.lastUsed = time.Now()
				atomic.AddInt64(&worker.processed, 1)
				atomic.AddInt64(&p.metrics.queuedJobs, -1)

				if err != nil {
					atomic.AddInt64(&worker.errors, 1)
					atomic.AddInt64(&p.metrics.failedJobs, 1)
					p.results <- err
				} else {
					atomic.AddInt64(&p.metrics.completedJobs, 1)
				}

				// Update average processing time
				p.updateAvgProcessTime(time.Since(start))

			case <-p.ctx.Done():
				return

			case <-worker.done:
				return
			}
		}
	}()
}

func (p *Pool) manageScale() {
	ticker := time.NewTicker(p.config.ScaleInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.adjustWorkerCount()
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Pool) adjustWorkerCount() {
	p.mu.RLock()
	currentWorkers := len(p.workers)
	queueSize := len(p.jobs) + len(p.priorities)
	p.mu.RUnlock()

	switch {
	case queueSize > currentWorkers && currentWorkers < p.config.MaxWorkers:
		// Scale up
		p.startWorker()
	case queueSize < currentWorkers/2 && currentWorkers > p.config.MinWorkers:
		// Scale down
		p.removeIdleWorker()
	}
}

func (p *Pool) removeIdleWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	var oldestIdle *workerInfo
	oldestIdleTime := time.Now()

	for _, worker := range p.workers {
		if worker.lastUsed.Before(oldestIdleTime) {
			oldestIdle = worker
			oldestIdleTime = worker.lastUsed
		}
	}

	if oldestIdle != nil && time.Since(oldestIdle.lastUsed) > p.config.IdleTimeout {
		close(oldestIdle.done)
		delete(p.workers, oldestIdle.id)
	}
}

// Stop gracefully shuts down the worker pool
func (p *Pool) Stop() {
	p.cancelFunc()
	p.workersWg.Wait()
	close(p.results)
}

// Results returns a channel for receiving job results
func (p *Pool) Results() <-chan error {
	return p.results
}

// GetMetrics returns current pool metrics
func (p *Pool) GetMetrics() poolMetrics {
	return poolMetrics{
		activeWorkers:  atomic.LoadInt64(&p.metrics.activeWorkers),
		completedJobs:  atomic.LoadInt64(&p.metrics.completedJobs),
		failedJobs:     atomic.LoadInt64(&p.metrics.failedJobs),
		queuedJobs:     atomic.LoadInt64(&p.metrics.queuedJobs),
		avgProcessTime: time.Duration(atomic.LoadInt64((*int64)(&p.metrics.avgProcessTime))),
	}
}

func (p *Pool) updateAvgProcessTime(newDuration time.Duration) {
	for {
		current := atomic.LoadInt64((*int64)(&p.metrics.avgProcessTime))
		completed := atomic.LoadInt64(&p.metrics.completedJobs)
		newAvg := time.Duration((int64(current)*completed + int64(newDuration)) / (completed + 1))

		if atomic.CompareAndSwapInt64((*int64)(&p.metrics.avgProcessTime), current, int64(newAvg)) {
			break
		}
	}
}
