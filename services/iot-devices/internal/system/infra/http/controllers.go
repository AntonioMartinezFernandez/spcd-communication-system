package system_http

import (
	"net/http"

	system_application "github.com/AntonioMartinezFernandez/services/iot-devices/internal/system/application"

	amf_query_bus "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus/query"
	amf_json_api "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api"
	json_api_response "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api/response"
)

func NewGetHealthcheckController(
	queryBus *amf_query_bus.QueryBus,
	jarm *amf_json_api.JsonApiResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := system_application.NewGetHealthcheckQuery()
		queryResponse, err := queryBus.Ask(r.Context(), query)

		switch err.(type) {
		case nil:
			ctx, writer := r.Context(), w
			response := queryResponse.(system_application.GetHealthcheckQueryHandlerResponse)
			jarm.WriteResponse(ctx, writer, &response, http.StatusOK)
			return
		default:
			ctx, writer, errResponse := r.Context(), w, json_api_response.NewInternalServerErrorWithDetails(err.Error())
			jarm.WriteErrorResponse(ctx, writer, errResponse, http.StatusInternalServerError, err)
			return
		}
	}
}
