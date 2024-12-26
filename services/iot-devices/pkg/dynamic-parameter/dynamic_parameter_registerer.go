package dynamic_parameter

import (
	"net/http"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus/command"
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus/query"
	http_server "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/http-server"
	json_api "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api"
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"
)

type DynamicParameterRouteRegisterer func(router *http_server.Router)

type DynamicParameterRouterRegistererOpsFunc func(ops *DynamicParameterRouterRegistererOps)

type DynamicParameterRouterRegistererOps struct {
	FetchPath  string
	UpdatePath string

	CommandBus command.Bus
	QueryBus   query.Bus
	Logger     logger.Logger

	AuthMiddleware http_server.Middleware
}

func (rro *DynamicParameterRouterRegistererOps) apply(ops ...DynamicParameterRouterRegistererOpsFunc) {
	for _, op := range ops {
		op(rro)
	}
}

func NewDefaultDynamicParameterRouterRegistererOps() *DynamicParameterRouterRegistererOps {
	return &DynamicParameterRouterRegistererOps{
		FetchPath:  "/system/dynamic-parameters/{parameterName}",
		UpdatePath: "/system/dynamic-parameters/{parameterName}",

		CommandBus: nil,
		QueryBus:   nil,
		Logger:     nil,

		AuthMiddleware: defaultAuthMiddleware(),
	}
}

func defaultAuthMiddleware() http_server.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func WithFetchPath(path string) DynamicParameterRouterRegistererOpsFunc {
	return func(ops *DynamicParameterRouterRegistererOps) {
		ops.FetchPath = path
	}
}

func WithUpdatePath(path string) DynamicParameterRouterRegistererOpsFunc {
	return func(ops *DynamicParameterRouterRegistererOps) {
		ops.UpdatePath = path
	}
}

func WithAuthMiddleware(middleware http_server.Middleware) DynamicParameterRouterRegistererOpsFunc {
	return func(ops *DynamicParameterRouterRegistererOps) {
		ops.AuthMiddleware = middleware
	}
}

func WithCommandBus(commandBus command.Bus) DynamicParameterRouterRegistererOpsFunc {
	return func(ops *DynamicParameterRouterRegistererOps) {
		ops.CommandBus = commandBus
	}
}

func WithQueryBus(queryBus query.Bus) DynamicParameterRouterRegistererOpsFunc {
	return func(ops *DynamicParameterRouterRegistererOps) {
		ops.QueryBus = queryBus
	}
}

func WithLogger(logger logger.Logger) DynamicParameterRouterRegistererOpsFunc {
	return func(ops *DynamicParameterRouterRegistererOps) {
		ops.Logger = logger
	}
}

func RegisterDynamicParametersRoutes(ops ...DynamicParameterRouterRegistererOpsFunc) DynamicParameterRouteRegisterer {
	options := NewDefaultDynamicParameterRouterRegistererOps()
	options.apply(ops...)

	return func(router *http_server.Router) {
		responseMiddleware := json_api.NewJsonApiResponseMiddleware(options.Logger)

		router.Get(
			options.FetchPath,
			HandleGetDynamicParameter(options.QueryBus, responseMiddleware),
			options.AuthMiddleware,
		)

		router.Put(
			options.UpdatePath,
			HandleChangeDynamicParameter(options.CommandBus, responseMiddleware),
			options.AuthMiddleware,
		)
	}
}

func RegisterDynamicParameterBusesOperations(
	ulidProvider utils.UlidProvider,
	retriever *DynamicParameterRetriever,
	repository DynamicParameterRepository,
	commandBus command.Bus,
	queryBus query.Bus,
) {
	findQueryHandler := NewFindDynamicParameterQueryHandler(ulidProvider, retriever)
	changeCommandHandler := NewChangeDynamicParameterCommandHandler(retriever, repository)

	if err := queryBus.RegisterQuery(&FindDynamicParameterQuery{}, findQueryHandler); err != nil {
		panic(err)
	}

	if err := commandBus.RegisterCommand(&ChangeDynamicParameterCommand{}, changeCommandHandler); err != nil {
		panic(err)
	}
}
