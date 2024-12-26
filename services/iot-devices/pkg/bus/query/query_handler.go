package query

import (
	"context"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus"
)

type QueryHandler interface {
	Handle(ctx context.Context, query bus.Dto) (interface{}, error)
}
