package system_application_test

import (
	"context"
	"testing"

	system_application "github.com/AntonioMartinezFernandez/services/iot-devices/internal/system/application"
	system_domain_mocks "github.com/AntonioMartinezFernandez/services/iot-devices/internal/system/domain/mocks"

	amf_utils "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetHealthcheckQueryHandler(t *testing.T) {
	t.Run("should return system ok", func(t *testing.T) {
		// Test data
		ctx := context.Background()
		serviceName := "service-name"
		statuses := map[string]string{"system": "ok"}

		// Test Dependencies
		ulidProvider := amf_utils.NewFixedUlidProvider()
		ulid := ulidProvider.New().String()

		healthchecker := system_domain_mocks.NewHealthChecker(t)
		healthchecker.On("Check", ctx).Return(statuses, nil).Once()

		// SUT
		handler := system_application.NewGetHealthcheckQueryHandler(serviceName, ulidProvider, healthchecker)

		// Test execution and assertions
		query := system_application.NewGetHealthcheckQuery()
		queryResponse, err := handler.Handle(ctx, query)
		assert.NoError(t, err)
		castedQueryResponse := queryResponse.(system_application.GetHealthcheckQueryHandlerResponse)
		assert.Equal(t, castedQueryResponse.Id, ulid)
		assert.Equal(t, castedQueryResponse.ServiceName, serviceName)
		assert.Equal(t, castedQueryResponse.Status, statuses)
	})
}
