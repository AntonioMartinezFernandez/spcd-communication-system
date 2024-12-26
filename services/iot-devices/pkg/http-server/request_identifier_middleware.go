package http_server

import (
	"context"
	"net/http"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"
)

type correlationId string

const HeaderRequestIdentifier = "X-Request-Id"
const contextKeyRequestIdentifier correlationId = "correlation_id"

type IdentifierGenerator func() string

type RequestIdentifierMiddleware struct {
	idGenerator IdentifierGenerator
}

func NewRequestIdentifierMiddleware(generator IdentifierGenerator) *RequestIdentifierMiddleware {
	return &RequestIdentifierMiddleware{idGenerator: generator}
}

func (rim *RequestIdentifierMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rim.idGenerator == nil {
			rim.idGenerator = rim.defaultIdentifierGenerator()
		}

		requestId := rim.fetchOrGenerateIdentifier(r)

		r.Header.Set(HeaderRequestIdentifier, requestId)

		newRequest := r.WithContext(context.WithValue(r.Context(), contextKeyRequestIdentifier, requestId))

		w.Header().Set(HeaderRequestIdentifier, requestId)

		next.ServeHTTP(w, newRequest)
	})
}

func (rim *RequestIdentifierMiddleware) fetchOrGenerateIdentifier(r *http.Request) string {
	requestId := r.Header.Get(HeaderRequestIdentifier)

	if requestId == "" {
		requestId = rim.idGenerator()
	}

	return requestId
}

func (rim *RequestIdentifierMiddleware) defaultIdentifierGenerator() IdentifierGenerator {
	return func() string {
		return utils.NewRandomUuidProvider().New().String()
	}
}
