package hawk_workflow

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type WorkflowCreateParams struct {
	Namespace  string
	MethodName func(job *work.Job)
	Cron       string
	RedisPort  string
}

type Workflow struct {
	redisPool  *redis.Pool
	enqueuer   *work.Enqueuer
	workerPool *work.WorkerPool
	namespace  string
	methodName func(job *work.Job)
	cron       string
}

type Context struct {
}

func CreateNewWorkflow(params WorkflowCreateParams) *Workflow {
	w := Workflow{}
	w.namespace = params.Namespace
	w = w.createRedisPool(params.RedisPort)
	w = w.createEnqueuer()
	w.methodName = params.MethodName
	w = w.createWorkerPool()
	w.cron = params.Cron
	return &w
}

func (w Workflow) createRedisPool(port string) Workflow {
	// You need an active redis connection
	w.redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf(":%s", port))
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
	methodNameString := fmt.Sprintf("%v", w.methodName)
	w.workerPool.PeriodicallyEnqueue(w.cron, methodNameString)
	w.workerPool.Job(methodNameString, w.methodName)
	w.workerPool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	w.workerPool.Stop()
}
