package dynamic_parameter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dynamic_parameter "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/dynamic-parameter"
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"
)

func TestFindDynamicParameter(t *testing.T) {
	rootCtx := context.Background()
	repository := new(DynamicParameterRepositoryMock)
	parameters := map[string]interface{}{"test_flag": true}
	ulidProvider := utils.NewFixedUlidProvider()
	handler := dynamic_parameter.NewFindDynamicParameterQueryHandler(
		ulidProvider,
		dynamic_parameter.NewDynamicParameterRetriever(repository, parameters),
	)

	t.Run("Find dynamic parameter without errors", func(t *testing.T) {
		query := &dynamic_parameter.FindDynamicParameterQuery{Name: "test_flag"}

		value := parameters[query.Name]
		parameterResponse := &dynamic_parameter.DynamicParameterResponse{
			ID:           ulidProvider.New().String(),
			Name:         query.Name,
			DefaultValue: value,
			DynamicValue: &value,
		}

		repository.ShouldSearchDynamicParameter(rootCtx, dynamic_parameter.ParameterName(query.Name), &value)
		response, err := handler.Handle(rootCtx, query)

		assert.NoError(t, err)
		assert.Equal(t, response, parameterResponse)

		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Find dynamic when dynamic values is not stored", func(t *testing.T) {
		query := &dynamic_parameter.FindDynamicParameterQuery{Name: "test_flag"}

		parameterResponse := &dynamic_parameter.DynamicParameterResponse{
			ID:           ulidProvider.New().String(),
			Name:         query.Name,
			DefaultValue: parameters[query.Name],
			DynamicValue: nil,
		}

		repository.ShouldSearchDynamicParameter(rootCtx, dynamic_parameter.ParameterName(query.Name), nil)
		response, err := handler.Handle(rootCtx, query)

		assert.NoError(t, err)
		assert.Equal(t, response, parameterResponse)

		mock.AssertExpectationsForObjects(t, repository)
	})
}

func TestFindDynamicParameterFail(t *testing.T) {
	rootCtx := context.Background()
	repository := new(DynamicParameterRepositoryMock)
	parameters := map[string]interface{}{"test_flag": true}
	ulidProvider := utils.NewFixedUlidProvider()
	handler := dynamic_parameter.NewFindDynamicParameterQueryHandler(
		ulidProvider,
		dynamic_parameter.NewDynamicParameterRetriever(repository, parameters),
	)

	t.Run("Dynamic parameter is not mapped in configuration", func(t *testing.T) {
		query := &dynamic_parameter.FindDynamicParameterQuery{Name: "invalid_param"}

		_, err := handler.Handle(rootCtx, query)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Dynamic parameter not exists")

		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Error getting dynamic parameter from repository", func(t *testing.T) {
		query := &dynamic_parameter.FindDynamicParameterQuery{Name: "test_flag"}

		repository.ShouldSearchDynamicParameterAndFail(rootCtx, dynamic_parameter.ParameterName(query.Name), nil, errors.New("some error"))
		_, err := handler.Handle(rootCtx, query)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "some error")

		mock.AssertExpectationsForObjects(t, repository)
	})
}
