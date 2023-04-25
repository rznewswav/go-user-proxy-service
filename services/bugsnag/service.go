package bugsnag

import (
	"errors"
	"fmt"
	"log"
	"os"
	"service/services/common/utils"
	"service/services/stack"

	bugsnag_interfaces "service/services/bugsnag/interfaces"
	bugsnag_structs "service/services/bugsnag/structs"

	"github.com/bugsnag/bugsnag-go/v2"
)

func IsDeploymenReleaseStage(stage string) bool {
	return stage == "staging" || stage == "production"
}

func DefaultReleaseStages() func() []string {
	return func() []string {
		return []string{"staging", "production"}
	}
}

func ExtendReleaseStages(stage string) func() []string {
	if IsDeploymenReleaseStage(stage) {
		return DefaultReleaseStages()
	}
	return func() []string {
		return []string{"staging", "production", stage}
	}
}

var configured bool = false

/*
*
Only release stage 'staging' and 'production' will be enabled.
However, if `forceReleaseStage` is true, the unrelated release
stage will be enabled. This is useful to debug error reported
to bugsnag in development mode.
*/
func Start(
	apiKey string,
	releaseStage string,
	forceReleaseStage bool,
) {
	// let's just collect the last 8 events
	breadcrumbs.Init(8)
	logger := log.New(os.Stdout, "BUGSNAG ", log.LstdFlags)
	if configured {
		logger.Printf(
			"WARN: bugsnag is already initialized for stage: %s\n",
			releaseStage,
		)
		return
	}

	if !forceReleaseStage &&
		!IsDeploymenReleaseStage(releaseStage) {
		logger.Printf(
			"not starting bugsnag for stage: %s\n",
			releaseStage,
		)
		return
	} else {
		logger.Printf("initialized bugsnag for stage: %s\n", releaseStage)
	}

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       apiKey,
		ReleaseStage: releaseStage,
		NotifyReleaseStages: utils.Iff(
			forceReleaseStage,
			ExtendReleaseStages(releaseStage),
			DefaultReleaseStages(),
		),
		AppType: "worker",
		// ProjectPackages:     []string{"main", "service/**"},
	})

	bugsnag.OnBeforeNotify(
		func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			var payload *bugsnag_structs.NotifyableError

			switch {
			case errors.As(event.Error.Err, &payload):
				event.ErrorClass = payload.Class
				event.Stacktrace = payload.Stacks
				event.MetaData.AddStruct(
					"cause",
					payload.CausedBy,
				)
				for _, decorator := range payload.ErrorDecorators {
					decorator.Decorate(event, config)
				}
			}

			for i, breadcrumb := range GetBreadcrumbs() {
				event.MetaData.Add(
					"breadcrumb",
					fmt.Sprintf("%d", i),
					breadcrumb,
				)

			}

			return nil
		},
	)

	configured = true
}

func GetHandler() bugsnag_interfaces.IBugsnagHandler {
	if configured {
		return &bugsnag_structs.BugsnagHandler{}
	}
	return &bugsnag_structs.DevBugsnagHandler{}
}

type ErrorClass string

func New(
	ErrorClass ErrorClass,
) *bugsnag_structs.NotifyableError {
	ne := new(bugsnag_structs.NotifyableError)
	ne.Stacks = stack.GetStackTrace()

	ne.Class = string(ErrorClass)
	return ne
}

func FromError(
	ErrorClass string,
	sourceError error,
) *bugsnag_structs.NotifyableError {
	ne := new(bugsnag_structs.NotifyableError)

	// inherit from parent
	if notifiableError, ok := sourceError.(*bugsnag_structs.NotifyableError); ok {
		ne.Stacks = notifiableError.Stacks
		ne.ErrorDecorators = notifiableError.ErrorDecorators
	} else {
		ne.Stacks = stack.GetStackTrace()
	}

	ne.Class = ErrorClass
	ne.SetCausedBy(sourceError)
	return ne
}
