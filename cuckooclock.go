package hawk_workflow

import (
	"fmt"

	"github.com/gocraft/work"
)

// Any method you create needs this method signature:
// func MethodName(job *work.Job) error
func RunCuckooClock(job *work.Job) error {
	fmt.Println("CUCKOO")
	return nil
}
