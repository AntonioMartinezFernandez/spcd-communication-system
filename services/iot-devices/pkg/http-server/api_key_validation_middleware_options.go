package http_server

import (
	"strings"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"
)

const defaultApiKeyHeaderName = "X-Api-Key"

type ApiKeyValidationMiddlewareOpsFunc func(*ApiKeyValidationMiddlewareOps)

type ApiKeyValidationMiddlewareOps struct {
	headerName string
	logger     logger.Logger
	keys       map[string]StaticApiKey
}

func NewDefaultApiKeyValidationMiddlewareOps() *ApiKeyValidationMiddlewareOps {
	return &ApiKeyValidationMiddlewareOps{
		headerName: defaultApiKeyHeaderName,
		logger:     nil,
		keys:       make(StaticApiKeyStorage),
	}
}

func WithHeaderName(header string) ApiKeyValidationMiddlewareOpsFunc {
	return func(ops *ApiKeyValidationMiddlewareOps) {
		ops.headerName = header
	}
}

func WithLogger(logger logger.Logger) ApiKeyValidationMiddlewareOpsFunc {
	return func(ops *ApiKeyValidationMiddlewareOps) {
		ops.logger = logger
	}
}

func WithDefaultKeyOwner(key string) ApiKeyValidationMiddlewareOpsFunc {
	return func(ops *ApiKeyValidationMiddlewareOps) {
		ops.keys = map[string]StaticApiKey{
			key: NewDefaultStaticApiKey(key),
		}
	}
}

func WithKeysByOwner(keys ...StaticApiKey) ApiKeyValidationMiddlewareOpsFunc {
	return func(ops *ApiKeyValidationMiddlewareOps) {
		storage := make(StaticApiKeyStorage, len(keys))
		for _, key := range keys {
			storage[key.Key] = key
		}

		ops.keys = storage
	}
}

// StaticApiKeysFromPipedString create a collection of identifiable api keys
// from string which follows the format <owner>,<key>|<owner>,<key>|<owner>,<key>
func StaticApiKeysFromPipedString(rawKeys string) []StaticApiKey {
	pipedKeys := strings.Split(strings.TrimSpace(rawKeys), "|")
	keys := make([]StaticApiKey, 0)

	for _, pipedKey := range pipedKeys {
		rawKey := strings.Split(strings.TrimSpace(pipedKey), ",")
		if len(rawKey) != 2 {
			continue
		}

		keyOwner, keyValue := rawKey[0], rawKey[1]

		keys = append(keys, NewStaticApiKey(keyOwner, keyValue))
	}

	return keys
}
