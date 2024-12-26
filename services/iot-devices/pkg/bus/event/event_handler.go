package event

import (
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus"
)

type EventHandler interface {
	Handle(event bus.Event) error
}
