package bus

type Dto interface {
	Type() string
}

type InvalidDto struct {
	message string
}

func NewInvalidDto(message string) *InvalidDto {
	return &InvalidDto{message: message}
}

func (i InvalidDto) Error() string {
	return i.message
}

type BlockOperationCommand interface {
	Dto
	BlockingKey() string
}
