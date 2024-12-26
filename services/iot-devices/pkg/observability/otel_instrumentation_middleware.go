package observability

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type OtelInstrumentationMiddleware struct {
	serviceName string
}

func NewOtelInstrumentationMiddleware(serviceName string) *OtelInstrumentationMiddleware {
	return &OtelInstrumentationMiddleware{serviceName: serviceName}
}

func (oim *OtelInstrumentationMiddleware) Middleware(next http.Handler) http.Handler {
	omm := otelmux.Middleware(oim.serviceName)
	return omm.Middleware(next)
}
