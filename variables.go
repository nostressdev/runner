package runner

type VariableProvider interface {
	GetString(group string, variable string) string
	GetInt(group string, variable string) int
}
