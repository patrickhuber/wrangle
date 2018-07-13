package renderers

type Renderer interface {
	RenderEnvironment(environmentVariables map[string]string) string
	RenderEnvironmentVariable(variable string, value string) string
	RenderProcess(path string, args []string, environmentVariables map[string]string) string
	Shell() string
}
