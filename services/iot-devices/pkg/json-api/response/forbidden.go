package json_api_response

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"
)

const (
	forbiddenDefaultTitle = "Forbidden"
	forbiddenDefaultCode  = "forbidden"
)

func NewForbidden(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   forbiddenDefaultCode,
		Title:  forbiddenDefaultTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusForbidden),
	}}
}

func NewForbiddenWithDetails(detail string, items ...MetadataItem) []*jsonapi.ErrorObject {
	metadata := NewMetadata(items...).MetadataMap()

	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   forbiddenDefaultCode,
		Title:  forbiddenDefaultTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusForbidden),
		Meta:   &metadata,
	}}
}
