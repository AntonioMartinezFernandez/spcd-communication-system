package http_server

import (
	"log/slog"
	"net/http"

	json_api "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api"
	json_api_response "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-api/response"
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"
)

const (
	defaultStaticApiKeyOwner            = "default"
	invalidKeyProvidedErrorMessage      = "invalid key provided"
	staticApiKeyUsageInformationMessage = "an static api key has been used"
)

type StaticApiKey struct {
	Owner string
	Key   string
}

func NewStaticApiKey(owner, key string) StaticApiKey {
	return StaticApiKey{
		Owner: owner,
		Key:   key,
	}
}

func NewDefaultStaticApiKey(key string) StaticApiKey {
	return NewStaticApiKey(defaultStaticApiKeyOwner, key)
}

type StaticApiKeyStorage map[string]StaticApiKey

func (sks StaticApiKeyStorage) SearchByKey(key string) (StaticApiKey, bool) {
	if value, ok := sks[key]; ok {
		return value, true
	}

	return StaticApiKey{}, false
}

type ApiKeyValidationMiddleware struct {
	headerName string
	logger     logger.Logger
	keys       StaticApiKeyStorage

	responseMiddleware *json_api.JsonApiResponseMiddleware
}

func NewApiKeyValidationMiddleware(
	responseMiddleware *json_api.JsonApiResponseMiddleware,
	ops ...ApiKeyValidationMiddlewareOpsFunc,
) *ApiKeyValidationMiddleware {
	options := NewDefaultApiKeyValidationMiddlewareOps()
	for _, op := range ops {
		op(options)
	}

	return &ApiKeyValidationMiddleware{
		headerName: options.headerName,
		logger:     options.logger,
		keys:       options.keys,

		responseMiddleware: responseMiddleware,
	}
}

func (kvm *ApiKeyValidationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if apiKey := req.Header.Get(kvm.headerName); apiKey != "" {
			if key, exists := kvm.keys.SearchByKey(apiKey); exists {
				kvm.registerStaticApiKeyUsage(key, req)
				next.ServeHTTP(w, req)
				return
			}
		}

		kvm.writeUnauthorizedResponse(w, req)
	})
}

func (kvm *ApiKeyValidationMiddleware) writeUnauthorizedResponse(
	w http.ResponseWriter,
	req *http.Request,
) {
	if kvm.logger != nil {
		kvm.logger.Warn(req.Context(), invalidKeyProvidedErrorMessage)
	}

	err, code := json_api_response.NewUnauthorized(invalidKeyProvidedErrorMessage), http.StatusUnauthorized
	kvm.responseMiddleware.WriteErrorResponse(req.Context(), w, err, code, nil)
}

func (kvm *ApiKeyValidationMiddleware) registerStaticApiKeyUsage(
	key StaticApiKey,
	req *http.Request,
) {
	if kvm.logger != nil {
		kvm.logger.Info(
			req.Context(),
			staticApiKeyUsageInformationMessage,
			slog.String("static_key_owner", key.Owner),
			slog.String("path", req.URL.Path),
		)
	}
}
