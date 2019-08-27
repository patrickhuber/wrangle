package filesystem 

// NewOs creates a new os file system
func NewOs() FileSystem{
	return newAferoOsFileSystem()
}