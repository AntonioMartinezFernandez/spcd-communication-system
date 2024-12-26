package sqldb

import "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/domain"

const invalidPoolConfigProvidedErrorMessage = "Invalid pool config provided"

type InvalidPoolConfigProvided struct {
	domain.RootCriticalError
	items map[string]interface{}
}

func (icp *InvalidPoolConfigProvided) Error() string {
	return invalidPoolConfigProvidedErrorMessage
}

func (icp *InvalidPoolConfigProvided) ExtraItems() map[string]interface{} {
	return icp.items
}

func NewInvalidPoolConfigProvided(driverName string) *InvalidPoolConfigProvided {
	return &InvalidPoolConfigProvided{
		items: map[string]interface{}{
			"driver": driverName,
		},
	}
}
