package filesystem

import (
	"fmt"
)

// File Extensions
const (
	// ExtensionVmt Material file extension
	ExtensionVmt = ".vmt"
	// ExtensionVtf Texture file extension
	ExtensionVtf = ".vtf"
)

// FilePath prefixes
const (
	// BasePathMaterial is path prefix for all materials/textures
	BasePathMaterial = "materials/"
	// BasePathModels is path prefix for all models/props
	BasePathModels = "models/"
)

type FileNotFoundError struct {
	fileName string
}

func (err FileNotFoundError) Error() string {
	return fmt.Sprintf("%s not found in filesystem", err.fileName)
}

func NewFileNotFoundError(filename string) *FileNotFoundError {
	return &FileNotFoundError{
		fileName: filename,
	}
}
