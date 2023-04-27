package bugsnag

func IsDeploymentReleaseStage(stage string) bool {
	return stage == "staging" || stage == "production"
}

func DefaultReleaseStages() func() []string {
	return func() []string {
		return []string{"staging", "production"}
	}
}

func ExtendReleaseStages(stage string) func() []string {
	if IsDeploymentReleaseStage(stage) {
		return DefaultReleaseStages()
	}
	return func() []string {
		return []string{"staging", "production", stage}
	}
}
