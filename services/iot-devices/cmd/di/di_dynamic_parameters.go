package di

import (
	"fmt"

	fgu_dynamic_parameter "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/dynamic-parameter"
	fgu_http_server "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/http-server"
)

const changeDynamicParameterJsonSchemaFileName = "change-dynamic-parameter.schema.json"

type DynamicParameterServices struct {
	DynamicParameterRetriever *fgu_dynamic_parameter.DynamicParameterRetriever
	StaticApiKeyStorage       []fgu_http_server.StaticApiKey
}

func InitDynamicParameterServices(commonServices *CommonServices, httpServices *HttpServices) *DynamicParameterServices {
	repository := fgu_dynamic_parameter.NewRedisDynamicParameterRepository(commonServices.RedisClient)
	retriever := fgu_dynamic_parameter.NewDynamicParameterRetrieverFromConfigFile(repository, commonServices.Config.DynamicParametersFilePath)

	fgu_dynamic_parameter.RegisterDynamicParameterBusesOperations(
		commonServices.UlidProvider,
		retriever,
		repository,
		commonServices.CommandBus,
		commonServices.QueryBus,
	)

	staticApiKeys := fgu_http_server.StaticApiKeysFromPipedString(commonServices.Config.DynamicParametersApiKeys)

	dynamicParametersServices := &DynamicParameterServices{
		DynamicParameterRetriever: retriever,
		StaticApiKeyStorage:       staticApiKeys,
	}

	registerDynamicParameterRoutes(dynamicParametersServices, commonServices, httpServices)

	return dynamicParametersServices
}

func registerDynamicParameterRoutes(
	dynamicParametersServices *DynamicParameterServices,
	commonServices *CommonServices,
	httpServices *HttpServices,
) {
	staticApiKeysMiddleware := fgu_http_server.NewApiKeyValidationMiddleware(
		httpServices.JsonApiResponseMiddleware,
		fgu_http_server.WithLogger(commonServices.Logger),
		fgu_http_server.WithKeysByOwner(dynamicParametersServices.StaticApiKeyStorage...),
	)

	changeDynamicParameterJsonSchemaValidator := fgu_http_server.NewRequestValidatorMiddleware(
		httpServices.JsonApiResponseMiddleware,
		fmt.Sprintf(
			"%s/%s/%s",
			commonServices.Config.JsonSchemaBasePath,
			"dynamic-parameters",
			changeDynamicParameterJsonSchemaFileName,
		),
	)

	httpServices.Router.Get(
		"/system/parameter/{parameterName}",
		fgu_dynamic_parameter.HandleGetDynamicParameter(
			commonServices.QueryBus,
			httpServices.JsonApiResponseMiddleware,
		),
		staticApiKeysMiddleware.Middleware,
	)

	httpServices.Router.Put(
		"/system/parameter/{parameterName}",
		fgu_dynamic_parameter.HandleChangeDynamicParameter(
			commonServices.CommandBus,
			httpServices.JsonApiResponseMiddleware,
		),
		staticApiKeysMiddleware.Middleware,
		changeDynamicParameterJsonSchemaValidator.Middleware,
	)
}
