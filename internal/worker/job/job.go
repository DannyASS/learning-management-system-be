package jobs

type Job interface {
	Process() error
}

var JobQueue chan Job
