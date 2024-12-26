package dynamic_parameter_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	dynamic_parameter "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/dynamic-parameter"
)

type DynamicParameterRepositoryMock struct {
	mock.Mock
}

func (dp *DynamicParameterRepositoryMock) Search(ctx context.Context, name dynamic_parameter.ParameterName) (interface{}, error) {
	args := dp.Called(ctx, name)

	return args[0], args.Error(1)
}

func (dp *DynamicParameterRepositoryMock) Save(ctx context.Context, parameter dynamic_parameter.DynamicParameter) error {
	args := dp.Called(ctx, parameter)

	return args.Error(0)
}

func (dp *DynamicParameterRepositoryMock) ShouldSearchDynamicParameter(ctx context.Context, name dynamic_parameter.ParameterName, value interface{}) {
	dp.
		On("Search", ctx, name).
		Once().
		Return(value, nil)
}

func (dp *DynamicParameterRepositoryMock) ShouldSearchDynamicParameterAndFail(ctx context.Context, name dynamic_parameter.ParameterName, value interface{}, err error) {
	dp.
		On("Search", ctx, name).
		Once().
		Return(value, err)
}

func (dp *DynamicParameterRepositoryMock) ShouldSave(ctx context.Context, parameter dynamic_parameter.DynamicParameter) {
	dp.
		On("Save", ctx, parameter).
		Once().
		Return(nil)
}

func (dp *DynamicParameterRepositoryMock) ShouldSaveAndFail(ctx context.Context, parameter dynamic_parameter.ParameterValue, err error) {
	dp.
		On("Save", ctx, parameter).
		Once().
		Return(err)
}
