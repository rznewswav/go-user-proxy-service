package bugsnag

import (
	"fmt"
	"service/services/stack"

	driver "github.com/bugsnag/bugsnag-go/v2"
)

type Breadcrumb struct {
	Time   string
	Data   string
	Stacks []driver.StackFrame
}

func (bc *Breadcrumb) String() string {
	return fmt.Sprintf(
		"%s: %s\n%s",
		bc.Time,
		bc.Data,
		stack.SimpleString(bc.Stacks),
	)
}
