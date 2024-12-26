package dynamic_parameter

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus/command"
	http_server "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/http-server"
	json_api "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api"
	json_api_response "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api/response"
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"
)

func HandleChangeDynamicParameter(bus command.Bus, responseMiddleware *json_api.JsonApiResponseMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestParams, err := http_server.AllParamsRequest(r)
		if err != nil {
			responseMiddleware.WriteErrorResponse(
				r.Context(),
				w,
				json_api_response.NewInternalServerError(),
				http.StatusInternalServerError,
				err,
			)
			return
		}

		parameterName := mux.Vars(r)["parameterName"]
		parameterValue := utils.GetInMapValueOrDefault([]string{"data", "attributes", "value"}, requestParams, nil)
		cmd := &ChangeDynamicParameterCommand{Name: parameterName, Value: parameterValue}

		err = bus.Dispatch(r.Context(), cmd)

		switch err.(type) {
		case nil:
			ctx, writer, statusCode := r.Context(), w, http.StatusNoContent
			responseMiddleware.WriteResponse(ctx, writer, nil, statusCode)
			return
		case *DynamicParameterNotExists:
			ctx, writer, response := r.Context(), w, json_api_response.NewNotFound(err.Error())
			responseMiddleware.WriteErrorResponse(ctx, writer, response, http.StatusNotFound, err)
			return
		default:
			ctx, writer, response := r.Context(), w, json_api_response.NewInternalServerErrorWithDetails(err.Error())
			responseMiddleware.WriteErrorResponse(ctx, writer, response, http.StatusInternalServerError, err)
			return
		}
	}
}
