package bugsnag_structs

import (
	"fmt"
	bugsnag_interfaces "service/services/bugsnag/interfaces"

	driver "github.com/bugsnag/bugsnag-go/v2"
)

type NotifyableError struct {
	Class             string
	Message           string
	CausedBy          error
	RemediationAction string
	Stacks            []driver.StackFrame
	ErrorDecorators   []bugsnag_interfaces.BugsnagDecorator
}

func (ne *NotifyableError) SetMessage(
	Message string,
) *NotifyableError {
	ne.Message = Message
	return ne
}

func (ne *NotifyableError) AddBugsnagDecorator(
	Decorator bugsnag_interfaces.BugsnagDecorator,
) *NotifyableError {
	ne.ErrorDecorators = append(ne.ErrorDecorators, Decorator)
	return ne
}

func (ne *NotifyableError) SetCausedBy(
	CausedBy error,
) *NotifyableError {
	ne.CausedBy = CausedBy
	return ne
}

func (ne *NotifyableError) SetRemediationAction(
	action string,
) *NotifyableError {
	ne.RemediationAction = action
	return ne
}

func (ne *NotifyableError) Error() string {
	if ne.CausedBy != nil {
		return fmt.Errorf("%s: %w", ne.Message, ne.CausedBy).
			Error()
	}
	return fmt.Errorf("%s", ne.Message).Error()

}
