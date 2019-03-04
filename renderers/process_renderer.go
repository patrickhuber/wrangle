package renderers

type ProcessRenderer interface {
	RenderProcess(path string, args []string, environmentVariables map[string]string) string
}
