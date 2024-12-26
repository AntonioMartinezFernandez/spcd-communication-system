package dynamic_parameter

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus/query"
	json_api "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api"
	json_api_response "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api/response"
)

func HandleGetDynamicParameter(bus query.Bus, responseMiddleware *json_api.JsonApiResponseMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parameterName := mux.Vars(r)["parameterName"]
		response, err := bus.Ask(r.Context(), &FindDynamicParameterQuery{Name: parameterName})

		switch err.(type) {
		case nil:
			ctx, writer := r.Context(), w
			responseMiddleware.WriteResponse(ctx, writer, response, http.StatusOK)
			return
		case *DynamicParameterNotExists:
			ctx, writer, errResponse := r.Context(), w, json_api_response.NewNotFound(err.Error())
			responseMiddleware.WriteErrorResponse(ctx, writer, errResponse, http.StatusNotFound, err)
			return
		default:
			ctx, writer, errResponse := r.Context(), w, json_api_response.NewInternalServerErrorWithDetails(err.Error())
			responseMiddleware.WriteErrorResponse(ctx, writer, errResponse, http.StatusInternalServerError, err)
			return
		}
	}
}
