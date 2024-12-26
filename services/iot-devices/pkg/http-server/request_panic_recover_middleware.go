package http_server

import (
	"encoding/json"
	"net/http"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"
)

type PanicRecoverMiddleware struct {
	logger logger.Logger
}

func NewPanicRecoverMiddleware(l logger.Logger) *PanicRecoverMiddleware {
	return &PanicRecoverMiddleware{logger: l}
}

func (prm *PanicRecoverMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				prm.logger.Error(r.Context(), "Unhandled Error")
				jsonBody, _ := json.Marshal(map[string]interface{}{
					"errors": []map[string]string{
						{
							"title": "Internal Server Error",
						},
					},
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if _, er := w.Write(jsonBody); er != nil {
					prm.logger.Error(r.Context(), "Could not write the response of the panic error")
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
