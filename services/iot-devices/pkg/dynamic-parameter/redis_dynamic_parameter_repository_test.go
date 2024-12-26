package dynamic_parameter_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"

	dynamic_parameter "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/dynamic-parameter"
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"
)

type RedisDynamicParameterRepositoryTestSuite struct {
	suite.Suite
	redisClient *redis.Client
	logger      logger.Logger
	miniRedis   *miniredis.Miniredis
	ctx         context.Context
}

// SetupSuite is a setup function run before the suite starts
func (suite *RedisDynamicParameterRepositoryTestSuite) SetupSuite() {
	suite.logger = &logger.NullLogger{}

	// Create a miniRedis instance for testing
	miniRedis, err := miniredis.Run()
	suite.Require().NoError(err)
	suite.miniRedis = miniRedis

	// Create a Redis client with the miniRedis URL
	client := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	suite.redisClient = client
}

func (suite *RedisDynamicParameterRepositoryTestSuite) TearDownSuite() {
	suite.miniRedis.Close()
	_ = suite.redisClient.Close()
}

func (suite *RedisDynamicParameterRepositoryTestSuite) SetupTest() {
	suite.ctx = context.Background()

	// Set a redis feature flags that will be checked during the tests
	_ = suite.miniRedis.Set("dynamic_parameter:fake_boolean_flag", "1")
	suite.miniRedis.SetTTL("dynamic_parameter:fake_boolean_flag", 0)
}

func (suite *RedisDynamicParameterRepositoryTestSuite) TestSearchDynamicParameterWithSuccess() {
	repository := dynamic_parameter.NewRedisDynamicParameterRepository(suite.redisClient)
	parameterName := dynamic_parameter.ParameterName("fake_boolean_flag")
	value, err := repository.Search(suite.ctx, parameterName)

	suite.NoError(err)
	suite.Equal(float64(1), value)

	keys := suite.redisClient.Keys(suite.ctx, "dynamic_parameter:fake_boolean_flag").Val()
	suite.Equal(1, len(keys))
}

func (suite *RedisDynamicParameterRepositoryTestSuite) TestSearchDynamicParameterWillReturnNilIfNotExists() {
	repository := dynamic_parameter.NewRedisDynamicParameterRepository(suite.redisClient)
	parameterName := dynamic_parameter.ParameterName("not_existent_flag")
	value, err := repository.Search(suite.ctx, parameterName)

	suite.NoError(err)
	suite.Nil(value)
}

func (suite *RedisDynamicParameterRepositoryTestSuite) TestSaveDynamicParameterWithSuccess() {
	repository := dynamic_parameter.NewRedisDynamicParameterRepository(suite.redisClient)
	parameterName := dynamic_parameter.ParameterName("new_fancy_flag")

	err := repository.Save(suite.ctx, dynamic_parameter.DynamicParameter{
		Name:         parameterName,
		DefaultValue: false,
		DynamicValue: true,
	})

	suite.NoError(err)

	keys := suite.redisClient.Keys(suite.ctx, "dynamic_parameter:new_fancy_flag").Val()
	suite.Equal(1, len(keys))

	value := suite.redisClient.Get(suite.ctx, "dynamic_parameter:new_fancy_flag").Val()
	suite.Equal("true", value)
}

func TestDynamicParameterRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RedisDynamicParameterRepositoryTestSuite))
}
