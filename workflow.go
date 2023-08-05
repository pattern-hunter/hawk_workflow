package hawk_workflow

import (
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type WorkflowCreateParams struct {
	Namespace  string
	MethodName string
	Cron       string
}

type Workflow struct {
	redisPool  *redis.Pool
	enqueuer   *work.Enqueuer
	workerPool *work.WorkerPool
	namespace  string
	methodName string
	cron       string
}

type Context struct {
}

func CreateNewWorkflow(params WorkflowCreateParams) *Workflow {
	w := Workflow{}
	w.namespace = params.Namespace
	w = w.createRedisPool()
	w = w.createEnqueuer()
	w.methodName = params.MethodName
	w = w.createWorkerPool()
	w.cron = params.Cron
	return &w
}

func (w Workflow) createRedisPool() Workflow {
	// You need an active redis connection
	w.redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}
	return w
}

func (w Workflow) createEnqueuer() Workflow {
	w.enqueuer = work.NewEnqueuer(w.namespace, w.redisPool)
	return w
}

func (w Workflow) createWorkerPool() Workflow {
	w.workerPool = work.NewWorkerPool(Context{}, 10, w.namespace, w.redisPool)
	return w
}

func (w Workflow) RunPeriodicWorkflow() {
	methodMapWithoutContext := GenerateMethodMapWithoutContext()
	w.workerPool.PeriodicallyEnqueue(w.cron, w.methodName)
	w.workerPool.Job(w.methodName, methodMapWithoutContext[w.methodName])
	w.workerPool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	w.workerPool.Stop()
}
