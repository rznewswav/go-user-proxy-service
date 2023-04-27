package bugsnag

import (
	"errors"
	"fmt"
	"service/services/common/utils"
	"service/services/config"
	"service/services/stack"

	"github.com/bugsnag/bugsnag-go/v2"
)

var configured = false
var c = Config{}

func init() {
	c := config.QuietBuild(Config{})
	if c.NoReport {
		return
	}
}

// Start Only release stage 'staging' and 'production' will be enabled.
//
// However, if `forceReleaseStage` is true, the unrelated release
// stage will be enabled.
//
// This is useful to debug error reported
// to bugsnag in development mode.
func Start(
	logFn func(string, ...any),
) {
	apiKey := c.BugsnagApiKey
	releaseStage := c.AppEnv
	forceReleaseStage := c.Debug

	if c.NoReport {
		return
	}
	// let's just collect the last 8 events
	breadcrumbs.Init(8)
	if configured {
		logFn(
			"WARN: bugsnag is already initialized for stage: %s",
			releaseStage,
		)
		return
	}

	if !forceReleaseStage &&
		!IsDeploymentReleaseStage(releaseStage) {
		logFn(
			"not starting bugsnag for stage: %s",
			releaseStage,
		)
		return
	} else {
		logFn("initialized bugsnag for stage: %s", releaseStage)
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
	})

	bugsnag.OnBeforeNotify(
		func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			var payload *NotifiableError

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

func GetHandler() Handler {
	if configured {
		return &handler{}
	}
	return &silentHandler{}
}

type ErrorClass string

func New(
	ErrorClass ErrorClass,
) *NotifiableError {
	ne := new(NotifiableError)
	ne.Stacks = stack.GetStackTrace()

	ne.Class = string(ErrorClass)
	return ne
}

func FromError(
	ErrorClass string,
	sourceError error,
) *NotifiableError {
	ne := new(NotifiableError)

	// inherit from parent
	if notifiableError, ok := sourceError.(*NotifiableError); ok {
		ne.Stacks = notifiableError.Stacks
		ne.ErrorDecorators = notifiableError.ErrorDecorators
	} else {
		ne.Stacks = stack.GetStackTrace()
	}

	ne.Class = ErrorClass
	ne.SetCausedBy(sourceError)
	return ne
}
