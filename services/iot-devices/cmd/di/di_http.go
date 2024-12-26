package di

import (
	"github.com/AntonioMartinezFernandez/services/iot-devices/configs"

	amf_http_server "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/http-server"
	amf_json_api "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api"
	amf_observability "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/observability"
)

type HttpServices struct {
	Router                    *amf_http_server.Router
	JsonApiResponseMiddleware *amf_json_api.JsonApiResponseMiddleware
}

func InitHttpServices(commonServices *CommonServices) *HttpServices {
	return &HttpServices{
		Router:                    newRouter(commonServices.Config, commonServices),
		JsonApiResponseMiddleware: amf_json_api.NewJsonApiResponseMiddleware(commonServices.Logger),
	}
}

func newRouter(config configs.Config, commonServices *CommonServices) *amf_http_server.Router {
	return amf_http_server.DefaultRouter(
		config.HttpWriteTimeout,
		config.HttpReadTimeout,
		amf_http_server.NewPanicRecoverMiddleware(commonServices.Logger).Middleware,
		amf_observability.NewOtelInstrumentationMiddleware(commonServices.Config.AppServiceName).Middleware,
	)
}
