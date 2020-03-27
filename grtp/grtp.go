package grtp

/**
the pool hold inbound chan and job chan
*/
type Pool struct {
	InboundChannel chan *Task //inbound channel
	JobChannel     chan *Task //worker channel
	MaxWorkers     int        //go routine time
}
type PoolExecutor interface {
	Commit(t *Task) //commit task
	Run()           //go run
}

/**
syncPool:simple mode
*/
func New(maxWorker int) (*Pool, error) {
	if maxWorker < 1 {
		return nil, NewGrtpError("max worker would not less than 1.")
	}
	pool := Pool{
		InboundChannel: make(chan *Task),
		JobChannel:     make(chan *Task),
		MaxWorkers:     maxWorker,
	}
	return &pool, nil
}

/**
async Pool: this async is mean the taskQueue is async,can receive cap tasks.
			but the sync pool is mean that taskQueue is synchronized,only offer 1 task and poll one task.
*/
func NewAsync(maxWorker, asyncSize int) (*Pool, error) {
	if maxWorker < 1 || asyncSize < 1 {
		return nil, NewGrtpError("max worker or asyncQueueSize would not less than 1.")
	}
	pool := Pool{
		InboundChannel: make(chan *Task, asyncSize),
		JobChannel:     make(chan *Task, asyncSize),
		MaxWorkers:     maxWorker,
	}
	return &pool, nil
}

/**
commit task to pool
*/
func (p *Pool) Commit(t *Task) {
	//commit task
	p.InboundChannel <- t
}
func (p *Pool) Run() {
	//pre Start
	for i := 0; i < p.MaxWorkers; i++ {
		go p.worker(i)
	}
	//take inbound to job channel by current go routine.
	for task := range p.InboundChannel {
		p.JobChannel <- task
	}
	//finally
	defer close(p.JobChannel)
	defer close(p.InboundChannel)
}

/**
call way: go work(id)
*/
func (p *Pool) worker(id int) {
	//make a go routine always poll
	for task := range p.JobChannel {
		task.Execute()
	}
}
