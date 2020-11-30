package group

import (
	"time"
)

const (
	queueSize = 100
)

type Package struct {
	Host string
	From string
	To string
	Message []byte
}

type Queue struct {
	queue chan Package
	mutex <-chan time.Time
	timespan int
}

func (q *Queue) Init(timespan int) {
	q.queue = make(chan Package, queueSize)
	q.timespan = timespan
}

func (q *Queue) Enqueue(p Package) {
	q.queue <- p
}

func (q *Queue) Dequeue () Package{
	var res Package
	if q.mutex != nil {
		<-q.mutex
	}
	res = <- q.queue
	q.mutex = time.After(time.Duration(q.timespan)*time.Second)
	return res
}
