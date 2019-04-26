package settings

type Provider interface {
	Get() (*Settings, error)
	Set(s *Settings) error
}
