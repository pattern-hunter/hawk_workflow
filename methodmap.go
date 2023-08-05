package hawk_workflow

import (
	"github.com/gocraft/work"
)

func GenerateMethodMapWithoutContext() map[string]func(job *work.Job) error {
	methodMap := make(map[string]func(job *work.Job) error)

	// This is the place to add methods you want the workflow to run and map
	// them to more user-friendly names
	methodMap["cuckoo_clock"] = RunCuckooClock

	return methodMap
}
