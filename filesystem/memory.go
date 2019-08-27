package filesystem

// NewMemory creates a new in memory files system
func NewMemory() FileSystem {
	return newAferoMemoryFileSystem()
}
