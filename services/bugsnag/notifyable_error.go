package bugsnag

import (
	"fmt"
	driver "github.com/bugsnag/bugsnag-go/v2"
)

type NotifiableError struct {
	Class             string
	Message           string
	CausedBy          error
	RemediationAction string
	Stacks            []driver.StackFrame
	ErrorDecorators   []Decorator
}

func (ne *NotifiableError) SetMessage(
	Message string,
) *NotifiableError {
	ne.Message = Message
	return ne
}

func (ne *NotifiableError) AddBugsnagDecorator(
	Decorator Decorator,
) *NotifiableError {
	ne.ErrorDecorators = append(ne.ErrorDecorators, Decorator)
	return ne
}

func (ne *NotifiableError) SetCausedBy(
	CausedBy error,
) *NotifiableError {
	ne.CausedBy = CausedBy
	return ne
}

func (ne *NotifiableError) SetRemediationAction(
	action string,
) *NotifiableError {
	ne.RemediationAction = action
	return ne
}

func (ne *NotifiableError) Error() string {
	if ne.CausedBy != nil {
		return fmt.Errorf("%s: %w", ne.Message, ne.CausedBy).
			Error()
	}
	return fmt.Errorf("%s", ne.Message).Error()

}
