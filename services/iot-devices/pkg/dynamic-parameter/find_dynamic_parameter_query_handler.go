package dynamic_parameter

import (
	"context"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus"
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"
)

const findDynamicParameterQueryName = "find_dynamic_parameter_query"

type FindDynamicParameterQuery struct {
	Name string
}

func (fpq FindDynamicParameterQuery) Type() string {
	return findDynamicParameterQueryName
}

type FindDynamicParameterQueryHandler struct {
	ulidProvider utils.UlidProvider
	retriever    *DynamicParameterRetriever
}

func NewFindDynamicParameterQueryHandler(
	ulidProvider utils.UlidProvider,
	retriever *DynamicParameterRetriever,
) FindDynamicParameterQueryHandler {
	return FindDynamicParameterQueryHandler{ulidProvider: ulidProvider, retriever: retriever}
}

func (fd FindDynamicParameterQueryHandler) Handle(ctx context.Context, dto bus.Dto) (interface{}, error) {
	query, ok := dto.(*FindDynamicParameterQuery)
	if !ok {
		return nil, bus.NewInvalidDto("Invalid query")
	}

	parameter, err := fd.retriever.Get(ctx, ParameterName(query.Name))
	if err != nil {
		return nil, err
	}

	return NewDynamicParameterFromParameter(fd.ulidProvider.New().String(), parameter), nil
}
