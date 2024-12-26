package system_application

import (
	"context"

	system_domain "github.com/AntonioMartinezFernandez/services/iot-devices/internal/system/domain"

	amf_bus "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus"
	amf_utils "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"
)

type GetHealthcheckQueryHandler struct {
	serviceName   string
	ulidProvider  amf_utils.UlidProvider
	healthChecker system_domain.HealthChecker
}

type GetHealthcheckQueryHandlerResponse struct {
	Id          string            `json:"id" jsonapi:"primary,id"`
	Status      map[string]string `json:"status" jsonapi:"attr,status"`
	ServiceName string            `json:"serviceName" jsonapi:"attr,serviceName"`
}

func NewGetHealthcheckQueryHandler(
	serviceName string,
	ulidProvider amf_utils.UlidProvider,
	healthChecker system_domain.HealthChecker,
) *GetHealthcheckQueryHandler {
	return &GetHealthcheckQueryHandler{
		serviceName:   serviceName,
		ulidProvider:  ulidProvider,
		healthChecker: healthChecker,
	}
}

func (q GetHealthcheckQueryHandler) Handle(ctx context.Context, query amf_bus.Dto) (interface{}, error) {
	_, ok := query.(*GetHealthcheckQuery)
	if !ok {
		return nil, amf_bus.NewInvalidDto("invalid query")
	}

	statuses, err := q.healthChecker.Check(ctx)
	if err != nil {
		return nil, err
	}

	return GetHealthcheckQueryHandlerResponse{
		Id:          q.ulidProvider.New().String(),
		Status:      statuses,
		ServiceName: q.serviceName,
	}, nil
}
