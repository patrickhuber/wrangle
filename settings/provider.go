package settings

// Provider defines an interface for managing settings
type Provider interface {
	Get() (*Settings, error)
	Set(s *Settings) error
	Initialize() (*Settings, error)
}
