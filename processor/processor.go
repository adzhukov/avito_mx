package processor

import (
	"avito_mx/config"
	"avito_mx/models"
	"runtime"
	"sync"
)

type Processor struct {
	tasks        map[int64]models.Task
	queue        chan models.Task
	workersCount int
	wg           *sync.WaitGroup
}

func New() *Processor {
	p := Processor{
		tasks:        make(map[int64]models.Task),
		queue:        make(chan models.Task, 100),
		workersCount: runtime.NumCPU() * 2,
		wg:           &sync.WaitGroup{},
	}
	config.Queue = p.queue
	return &p
}

func (p *Processor) Start() {
	config.Logger.Printf("Starting %d worker pool\n", p.workersCount)
	for i := 0; i < p.workersCount; i++ {
		p.wg.Add(1)
		go p.Worker()
	}
}

func (p *Processor) Wait() {
	config.Logger.Println("Waiting until all workers exit")
	p.wg.Wait()
	config.Logger.Println("All workers exited")
}
