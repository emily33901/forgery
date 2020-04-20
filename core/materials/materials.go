package materials

import (
	"github.com/emily33901/forgery/core/textures"
	"github.com/golang-source-engine/vmt"
)

// Material
type Material struct {
	Props *vmt.Properties

	filePath string

	Textures struct {
		// Albedo
		Albedo *textures.Texture
		// Normal
		Normal *textures.Texture
	}
}

// Width returns this materials width. Albedo is used to
// determine material width where possible
func (mat *Material) Width() int {
	return mat.Textures.Albedo.Width()
}

// Height returns this materials height. Albedo is used to
// determine material height where possible
func (mat *Material) Height() int {
	return mat.Textures.Albedo.Height()
}

// FilePath returns this materials location in whatever
// filesystem it was found
func (mat *Material) FilePath() string {
	return mat.filePath
}

func (mat *Material) EvictTextures() {
	if mat.Textures.Albedo != nil {
		mat.Textures.Albedo.EvictFromMainMemory()
	}
	if mat.Textures.Normal != nil {
		mat.Textures.Normal.EvictFromMainMemory()
	}
}

func NewMaterial(filePath string, props *vmt.Properties) *Material {
	return &Material{
		filePath: filePath,
		Props:    props,
	}
}
