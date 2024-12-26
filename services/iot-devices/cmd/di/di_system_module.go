package di

import (
	system_application "github.com/AntonioMartinezFernandez/services/iot-devices/internal/system/application"
	system_infra "github.com/AntonioMartinezFernandez/services/iot-devices/internal/system/infra"
	system_http "github.com/AntonioMartinezFernandez/services/iot-devices/internal/system/infra/http"
)

type SystemServices struct {
	HealthcheckQueryHandler *system_application.GetHealthcheckQueryHandler
}

func InitSystemServices(commonServices *CommonServices, httpServices *HttpServices) *SystemServices {
	healthchecker := system_infra.NewSimpleHealthChecker(commonServices.Observability.Tracer)
	healthcheckQueryHandler := system_application.NewGetHealthcheckQueryHandler(
		commonServices.Config.AppServiceName,
		commonServices.UlidProvider,
		healthchecker,
	)

	systemServices := &SystemServices{
		HealthcheckQueryHandler: healthcheckQueryHandler,
	}

	registerSystemQueryHandlers(commonServices, systemServices)
	registerSystemRoutes(commonServices, httpServices)

	return systemServices
}

func registerSystemQueryHandlers(commonServices *CommonServices, systemServices *SystemServices) {
	registerQueryOrPanic(
		commonServices.QueryBus,
		&system_application.GetHealthcheckQuery{},
		systemServices.HealthcheckQueryHandler,
	)
}

func registerSystemRoutes(commonServices *CommonServices, httpServices *HttpServices) {
	httpServices.Router.Get(
		"/system/healthcheck",
		system_http.NewGetHealthcheckController(
			commonServices.QueryBus,
			httpServices.JsonApiResponseMiddleware,
		),
	)
}
