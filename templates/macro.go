package templates

// Macro is used to process a value
type Macro interface {
	Run(metadata *MacroMetadata) (string, error)
}
