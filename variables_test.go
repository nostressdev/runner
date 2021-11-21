package runner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MapProvider(t *testing.T) {
	m := NewMapVariableProvider(map[VariableKV]string{
		VariableKV{
			group:    "a",
			variable: "b",
		}: "xxx",
	})
	assert.NoError(t, m.EnsureEnvironmentVariables("a", "b"))
	assert.Error(t, m.EnsureEnvironmentVariables("a", "c"))
	assert.Error(t, m.EnsureEnvironmentVariables("b", "b"))
	assert.Error(t, m.EnsureEnvironmentVariables("b", "a"))
	assert.Equal(t, "xxx", m.GetString("a", "b"))
	assert.Equal(t, "", m.GetString("b", "a"))
}

func Test_MergeProviders(t *testing.T) {
	m1 := NewMapVariableProvider(map[VariableKV]string{
		VariableKV{
			group:    "a",
			variable: "b",
		}: "xxx",
	})
	m2 := NewMapVariableProvider(map[VariableKV]string{
		VariableKV{
			group:    "b",
			variable: "a",
		}: "yyy",
	})
	m := MergeProviders(m1, m2)
	assert.NoError(t, m.EnsureEnvironmentVariables("a", "b"))
	assert.Error(t, m.EnsureEnvironmentVariables("a", "c"))
	assert.Error(t, m.EnsureEnvironmentVariables("b", "b"))
	assert.NoError(t, m.EnsureEnvironmentVariables("b", "a"))
	assert.Equal(t, "xxx", m.GetString("a", "b"))
	assert.Equal(t, "yyy", m.GetString("b", "a"))
}
