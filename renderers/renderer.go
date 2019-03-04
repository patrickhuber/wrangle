package renderers

// Renderer defines a rendering interface
type Renderer interface {
	EnvironmentRenderer
	ProcessRenderer
	Format() string
}
