package dynamic_parameter

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

const prefix = "dynamic_parameter:"

type RedisDynamicParameterRepository struct {
	client *redis.Client
}

func NewRedisDynamicParameterRepository(client *redis.Client) *RedisDynamicParameterRepository {
	return &RedisDynamicParameterRepository{client: client}
}

func (r *RedisDynamicParameterRepository) Save(ctx context.Context, parameter DynamicParameter) error {
	bytes, err := json.Marshal(&parameter.DynamicValue)
	if err != nil {
		return err
	}
	_, err = r.client.Set(ctx, r.key(parameter.Name.Value()), bytes, 0).Result()

	return err
}

func (r *RedisDynamicParameterRepository) Search(ctx context.Context, name ParameterName) (interface{}, error) {
	var plainValue interface{}

	value, err := r.client.Get(ctx, r.key(name.Value())).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(value), &plainValue); err != nil {
		return nil, err
	}

	return plainValue, nil
}

func (r *RedisDynamicParameterRepository) key(value string) string {
	return prefix + value
}
