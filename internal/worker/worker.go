package worker

import (
	"log"
	"sync"
	"time"

	jobs "github.com/DannyAss/users/internal/worker/job"
)

var (
	wg       sync.WaitGroup
	quit     chan struct{}
	started  bool
	startMtx sync.Mutex
)

// InitJobQueue harus dipanggil sekali di bootstrap
func InitJobQueue(buffer int) {
	if jobs.JobQueue == nil {
		jobs.JobQueue = make(chan jobs.Job, buffer)
	}
	quit = make(chan struct{})
}

// StartWorkers menjalankan n worker
func StartWorkers(n int) {
	startMtx.Lock()
	if started {
		startMtx.Unlock()
		return
	}
	started = true
	startMtx.Unlock()

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go startWorker(i)
	}
}

func startWorker(id int) {
	defer wg.Done()
	log.Printf("[Worker %d] started", id)

	for {
		select {
		case <-quit:
			log.Printf("[Worker %d] received stop signal", id)
			return
		case job, ok := <-jobs.JobQueue:
			if !ok {
				// channel closed -> exit
				log.Printf("[Worker %d] job queue closed, exiting", id)
				return
			}

			// protect from panic di job.Process()
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[Worker %d] recovered panic: %v", id, r)
					}
				}()

				// eksekusi job (job.Process harus handle error sendiri)
				if err := job.Process(); err != nil {
					log.Printf("[Worker %d] job error: %v", id, err)
					// optional: retry, send to dead-letter queue, dll
				} else {
					log.Printf("[Worker %d] job finished", id)
				}
			}()
		}
	}
}

// StopWorkers memicu shutdown worker dan menunggu sampai selesai
func StopWorkers(timeout time.Duration) {
	// kirim sinyal stop
	close(quit)

	// optionally close job queue if you want
	// close(jobs.JobQueue)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("all workers stopped")
	case <-time.After(timeout):
		log.Println("stop workers timeout")
	}
}
