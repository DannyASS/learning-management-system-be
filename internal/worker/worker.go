package worker

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	jobs "github.com/DannyAss/users/internal/worker/job"
)

var (
	wg       sync.WaitGroup
	quit     chan struct{}
	started  bool
	startMtx sync.Mutex
)

// InitJobQueue harus dipanggil sekali di Buildapp
func InitJobQueue(buffer int) {
	if jobs.JobQueue == nil {
		jobs.JobQueue = make(chan jobs.Job, buffer)
	}
	quit = make(chan struct{})
}

// ---------------- Non-prefork workers ----------------
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
	log.Printf("[Worker %d] started, PID: %d", id, os.Getpid())

	for {
		select {
		case <-quit:
			log.Printf("[Worker %d] received stop signal", id)
			return
		case job, ok := <-jobs.JobQueue:
			if !ok {
				log.Printf("[Worker %d] job queue closed", id)
				return
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[Worker %d] recovered panic: %v", id, r)
					}
				}()

				if err := job.Process(); err != nil {
					log.Printf("[Worker %d] job error: %v", id, err)
				} else {
					log.Printf("[Worker %d] job finished", id)
				}
			}()
		}
	}
}

// ---------------- Prefork-safe workers ----------------
func StartWorkersPreforkSafe(n int, cfg *config.ConfigEnv) {
	for i := 0; i < n; i++ {
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Prefork Worker %d] recovered panic, PID: %d: %v", workerID, os.Getpid(), r)
				}
			}()

			log.Printf("[Prefork Worker %d] started, PID: %d", workerID, os.Getpid())

			// DB manager per worker
			db := database.NewDBManager(cfg.DBConnnect)
			if db == nil {
				log.Fatal("DB Manager init failed in worker", workerID)
			}

			// JobQueue per worker (fork-safe)
			jobQueue := make(chan jobs.Job, 1000)

			for {
				select {
				case <-quit:
					log.Printf("[Prefork Worker %d] received stop signal", workerID)
					return
				case job, ok := <-jobQueue:
					if !ok {
						log.Printf("[Prefork Worker %d] job queue closed", workerID)
						return
					}

					func() {
						defer func() {
							if r := recover(); r != nil {
								log.Printf("[Prefork Worker %d] recovered panic: %v", workerID, r)
							}
						}()

						if err := job.Process(); err != nil {
							log.Printf("[Prefork Worker %d] job error: %v", workerID, err)
						} else {
							log.Printf("[Prefork Worker %d] job finished", workerID)
						}
					}()
				}
			}
		}(i)
	}
}

// Stop workers
func StopWorkers(timeout time.Duration) {
	close(quit)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All workers stopped")
	case <-time.After(timeout):
		log.Println("Stop workers timeout")
	}
}
