package configs

import (
	"github.com/pkg/errors"
)

const (
	production  EnvironmentName = "production"
	test        EnvironmentName = "test"
	development EnvironmentName = "development"
)

type EnvironmentName string

type Environment struct {
	name EnvironmentName
}

func NewEnvironmentFromRawEnvVar(name string) Environment {
	guardEnvironmentName(name)
	return Environment{name: EnvironmentName(name)}
}

func environments() map[EnvironmentName]struct{} {
	return map[EnvironmentName]struct{}{
		production:  {},
		test:        {},
		development: {},
	}
}

func guardEnvironmentName(name string) {
	env := EnvironmentName(name)

	if _, environmentExists := environments()[env]; !environmentExists {
		panic(errors.Errorf("environment <%s> doesnt exist", name))
	}
}

func (env Environment) IsDevelopment() bool {
	return env.name == development || env.name == test
}

func (env Environment) IsTest() bool {
	return env.name == test
}
