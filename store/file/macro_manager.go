package file

type MacroManager interface {
	Run(metadata *MacroMetadata) (string, error)
	Register(name string, macro Macro)
}

type macroManager struct {
	macros map[string]Macro
}

func (manager *macroManager) Register(name string, macro Macro) {
	manager.macros[name] = macro
}

func (manager *macroManager) Run(metadata *MacroMetadata) (string, error) {
	macro := manager.macros[metadata.Name]
	return macro.Run(metadata)
}

func NewMacroManager() MacroManager {
	return &macroManager{
		macros: make(map[string]Macro),
	}
}
