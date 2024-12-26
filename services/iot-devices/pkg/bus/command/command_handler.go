package command

import (
	"context"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus"
)

type CommandHandler interface {
	Handle(ctx context.Context, command bus.Dto) error
}
