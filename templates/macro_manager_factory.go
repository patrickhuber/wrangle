package templates

type MacroManagerFactory interface {
	Create() MacroManager
}

type macroManagerFactory struct {
}

func (f *macroManagerFactory) Create() MacroManager {
	macroManager := NewMacroManager()
	macroManager.Register("ENC", NewEncryptionMacro())
	return macroManager
}

func NewMacroManagerFactory() MacroManagerFactory {
	return &macroManagerFactory{}
}
