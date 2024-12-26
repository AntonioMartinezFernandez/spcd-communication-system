package system_domain

import "context"

type HealthChecker interface {
	Check(ctx context.Context) (map[string]string, error)
}
