package dynamic_parameter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dynamic_parameter "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/dynamic-parameter"
)

func TestChangeDynamicParameter(t *testing.T) {
	rootCtx := context.Background()
	repository := new(DynamicParameterRepositoryMock)
	parameters := map[string]interface{}{"test_flag": "a value"}
	handler := dynamic_parameter.NewChangeDynamicParameterCommandHandler(dynamic_parameter.NewDynamicParameterRetriever(repository, parameters), repository)

	t.Run("Change dynamic without any error", func(t *testing.T) {
		command := &dynamic_parameter.ChangeDynamicParameterCommand{Name: "test_flag", Value: "newValue"}

		value := parameters[command.Name]

		repository.ShouldSearchDynamicParameter(rootCtx, dynamic_parameter.ParameterName(command.Name), &value)
		repository.ShouldSave(rootCtx, dynamic_parameter.DynamicParameter{
			Name:         dynamic_parameter.ParameterName(command.Name),
			DefaultValue: value,
			DynamicValue: command.Value,
		})
		err := handler.Handle(rootCtx, command)

		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, repository)
	})
}

func TestChangeDynamicParameterFail(t *testing.T) {
	rootCtx := context.Background()
	repository := new(DynamicParameterRepositoryMock)
	parameters := map[string]interface{}{"test_flag": true}
	handler := dynamic_parameter.NewChangeDynamicParameterCommandHandler(dynamic_parameter.NewDynamicParameterRetriever(repository, parameters), repository)

	t.Run("Dynamic parameter not mapped in configuration", func(t *testing.T) {
		command := &dynamic_parameter.ChangeDynamicParameterCommand{Name: "invalid_param"}

		err := handler.Handle(rootCtx, command)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Dynamic parameter not exists")

		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Error getting dynamic parameter from repository", func(t *testing.T) {
		command := &dynamic_parameter.ChangeDynamicParameterCommand{Name: "test_flag"}

		repository.ShouldSearchDynamicParameterAndFail(rootCtx, dynamic_parameter.ParameterName(command.Name), nil, errors.New("some error"))
		err := handler.Handle(rootCtx, command)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "some error")

		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Error saving dynamic parameter", func(t *testing.T) {
		command := &dynamic_parameter.ChangeDynamicParameterCommand{Name: "test_flag", Value: "newValue"}

		value := parameters[command.Name]

		repository.ShouldSearchDynamicParameter(rootCtx, dynamic_parameter.ParameterName(command.Name), &value)
		repository.ShouldSaveAndFail(rootCtx, dynamic_parameter.DynamicParameter{
			Name:         dynamic_parameter.ParameterName(command.Name),
			DefaultValue: value,
			DynamicValue: command.Value,
		}, errors.New("some error"))
		err := handler.Handle(rootCtx, command)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "some error")

		mock.AssertExpectationsForObjects(t, repository)
	})
}
