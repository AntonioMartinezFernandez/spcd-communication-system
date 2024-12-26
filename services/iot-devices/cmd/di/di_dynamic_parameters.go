package di

import (
	"fmt"

	amf_dynamic_parameter "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/dynamic-parameter"
	amf_http_server "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/http-server"
)

const changeDynamicParameterJsonSchemaFileName = "change-dynamic-parameter.schema.json"

type DynamicParameterServices struct {
	DynamicParameterRetriever *amf_dynamic_parameter.DynamicParameterRetriever
	StaticApiKeyStorage       []amf_http_server.StaticApiKey
}

func InitDynamicParameterServices(commonServices *CommonServices, httpServices *HttpServices) *DynamicParameterServices {
	repository := amf_dynamic_parameter.NewRedisDynamicParameterRepository(commonServices.RedisClient)
	retriever := amf_dynamic_parameter.NewDynamicParameterRetrieverFromConfigFile(repository, commonServices.Config.DynamicParametersFilePath)

	amf_dynamic_parameter.RegisterDynamicParameterBusesOperations(
		commonServices.UlidProvider,
		retriever,
		repository,
		commonServices.CommandBus,
		commonServices.QueryBus,
	)

	staticApiKeys := amf_http_server.StaticApiKeysFromPipedString(commonServices.Config.DynamicParametersApiKeys)

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
	staticApiKeysMiddleware := amf_http_server.NewApiKeyValidationMiddleware(
		httpServices.JsonApiResponseMiddleware,
		amf_http_server.WithLogger(commonServices.Logger),
		amf_http_server.WithKeysByOwner(dynamicParametersServices.StaticApiKeyStorage...),
	)

	changeDynamicParameterJsonSchemaValidator := amf_http_server.NewRequestValidatorMiddleware(
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
		amf_dynamic_parameter.HandleGetDynamicParameter(
			commonServices.QueryBus,
			httpServices.JsonApiResponseMiddleware,
		),
		staticApiKeysMiddleware.Middleware,
	)

	httpServices.Router.Put(
		"/system/parameter/{parameterName}",
		amf_dynamic_parameter.HandleChangeDynamicParameter(
			commonServices.CommandBus,
			httpServices.JsonApiResponseMiddleware,
		),
		staticApiKeysMiddleware.Middleware,
		changeDynamicParameterJsonSchemaValidator.Middleware,
	)
}
