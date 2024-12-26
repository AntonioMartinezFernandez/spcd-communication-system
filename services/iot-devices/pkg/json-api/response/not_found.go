package json_api_response

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"
)

const (
	notFoundCode  = "not_found"
	notFoundTitle = "Not Found"
)

func NewNotFound(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   notFoundCode,
		Title:  notFoundTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusNotFound),
	}}
}

func NewNotFoundErrorWithDetails(detail string, items ...MetadataItem) []*jsonapi.ErrorObject {
	metadata := NewMetadata(items...).MetadataMap()

	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   notFoundCode,
		Title:  notFoundTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusNotFound),
		Meta:   &metadata,
	}}
}
