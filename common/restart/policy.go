package restart

import (
	"github.com/docker/docker/api/types/container"
)

func AlwaysPolicy() container.RestartPolicy {
	return container.RestartPolicy{Name: "always"}
}

func NonePolicy() container.RestartPolicy {
	return container.RestartPolicy{Name: ""}
}

func UnlessStoppedPolicy() container.RestartPolicy {
	return container.RestartPolicy{Name: "unless-stopped"}
}

func OnFailurePolicy(retry int) container.RestartPolicy {
	return container.RestartPolicy{Name: "on-failure", MaximumRetryCount: retry}
}
