package bugsnag

import (
	bugsnag_structs "service/services/bugsnag/structs"
	"service/services/common/structs"
	"service/services/stack"
	"time"
)

var breadcrumbs = new(structs.RingLinkedList)

func LeaveBreadcrumb(data string) {
	bc := new(bugsnag_structs.Breadcrumb)
	bc.Time = time.Now().Format(time.RFC3339)
	bc.Data = data
	bc.Stacks = stack.GetStackTrace()[1:]
	breadcrumbs.Put(bc)
}

func GetBreadcrumbs() []string {
	return breadcrumbs.CollectNonNilToString()
}
