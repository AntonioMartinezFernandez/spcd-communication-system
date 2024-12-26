package http_server

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	amf_json_api "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api"
	response "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api/response"

	"github.com/google/jsonapi"
	"github.com/xeipuuv/gojsonschema"
)

type RequestValidatorMiddleware struct {
	responseMiddleware *amf_json_api.JsonApiResponseMiddleware
	schemaFilePath     string
}

func NewRequestValidatorMiddleware(
	responseMiddleware *amf_json_api.JsonApiResponseMiddleware,
	schemaFilePath string,
) *RequestValidatorMiddleware {
	return &RequestValidatorMiddleware{
		responseMiddleware: responseMiddleware,
		schemaFilePath:     schemaFilePath,
	}
}

func (rvm *RequestValidatorMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		absPath, _ := filepath.Abs(rvm.schemaFilePath)
		schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", absPath))
		bodyBytes, err := io.ReadAll(CloneRequest(r).Body)
		if err != nil {
			validationErrors := response.NewBadRequestForInvalidPayload()
			rvm.responseMiddleware.WriteErrorResponse(r.Context(), w, validationErrors, http.StatusBadRequest, err)
			return
		}

		documentLoader := gojsonschema.NewBytesLoader(bodyBytes)

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			validationErrors := response.NewBadRequestForInvalidPayload()
			rvm.responseMiddleware.WriteErrorResponse(r.Context(), w, validationErrors, http.StatusBadRequest, err)
			return
		}

		if !result.Valid() {
			var validationErrors []*jsonapi.ErrorObject
			errors := result.Errors()
			for i := range errors {
				desc := errors[i]
				var details map[string]interface{} = desc.Details()
				validationErrors = response.NewInvalidPayloadCustom(desc.Type(), desc.Description(), desc.String(), details)
			}
			rvm.responseMiddleware.WriteErrorResponse(r.Context(), w, validationErrors, http.StatusBadRequest, nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
