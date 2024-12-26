package dynamic_parameter

import (
	"context"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus"
)

const changeDynamicParameterCmdName = "change_dynamic_parameter_command"

type ChangeDynamicParameterCommand struct {
	Name  string
	Value interface{}
}

func (cdp *ChangeDynamicParameterCommand) Type() string {
	return changeDynamicParameterCmdName
}

type ChangeDynamicParameterCommandHandler struct {
	retriever  *DynamicParameterRetriever
	repository DynamicParameterRepository
}

func NewChangeDynamicParameterCommandHandler(retriever *DynamicParameterRetriever, repository DynamicParameterRepository) ChangeDynamicParameterCommandHandler {
	return ChangeDynamicParameterCommandHandler{retriever: retriever, repository: repository}
}

func (fd ChangeDynamicParameterCommandHandler) Handle(ctx context.Context, command bus.Dto) error {
	dpCommand, ok := command.(*ChangeDynamicParameterCommand)
	if !ok {
		return bus.NewInvalidDto("Invalid command")
	}

	parameter, err := fd.retriever.Get(ctx, ParameterName(dpCommand.Name))
	if err != nil {
		return err
	}

	updatedParameter := parameter.WithNewValue(dpCommand.Value)

	return fd.repository.Save(ctx, updatedParameter)
}
