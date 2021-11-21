package runner

import (
	"fmt"
	"os"
)

type VariableProvider interface {
	GetString(group string, variable string) string
	EnsureEnvironmentVariables(group string, variables ...string) error
}

var ErrVariableNotFound = fmt.Errorf("variable not found")

type environmentVariableProvider struct {
}

func (e *environmentVariableProvider) EnsureEnvironmentVariables(group string, variables ...string) error {
	for _, variable := range variables {
		if os.Getenv(group+"_"+variable) == "" {
			return ErrVariableNotFound
		}
	}
	return nil
}

func (e *environmentVariableProvider) GetString(group string, variable string) string {
	return os.Getenv(group + "_" + variable)
}

func NewEnvironmentVariableProvider() VariableProvider {
	return &environmentVariableProvider{}
}

type mergeProvider struct {
	providers []VariableProvider
}

func (m *mergeProvider) GetString(group string, variable string) string {
	for _, provider := range m.providers {
		if provider.EnsureEnvironmentVariables(group, variable) == nil {
			return provider.GetString(group, variable)
		}
	}
	return ""
}

func (m *mergeProvider) EnsureEnvironmentVariables(group string, variables ...string) error {
	for _, variable := range variables {
		found := false
		for _, provider := range m.providers {
			if provider.EnsureEnvironmentVariables(group, variable) == nil {
				found = true
				break
			}
		}
		if !found {
			return ErrVariableNotFound
		}
	}
	return nil
}

func MergeProviders(providers ...VariableProvider) VariableProvider {
	return &mergeProvider{
		providers: providers,
	}
}

type VariableKV struct {
	group    string
	variable string
}

type mapVariableProvider struct {
	m map[VariableKV]string
}

func (m *mapVariableProvider) GetString(group string, variable string) string {
	return m.m[VariableKV{group: group, variable: variable}]
}

func (m *mapVariableProvider) EnsureEnvironmentVariables(group string, variables ...string) error {
	for _, variable := range variables {
		if _, ok := m.m[VariableKV{group: group, variable: variable}]; !ok {
			return ErrVariableNotFound
		}
	}
	return nil
}

func NewMapVariableProvider(m map[VariableKV]string) VariableProvider {
	return &mapVariableProvider{m: m}
}
