package textures

import (
	"fmt"

	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/forgery/core/manager"
)

var textureManager *manager.Manager = manager.NewManager("")

func Load(path string, fs *filesystem.Filesystem) (*Texture, error) {
	v := textureManager.Get(path)

	if v != nil {
		return v.(*Texture), nil
	}

	// Try normalise
	normalPath := normaliseTexturePath(path)

	v = textureManager.Get(normalPath)
	if v != nil {
		return v.(*Texture), nil
	}

	// attempt to load
	fmt.Println("Loading", normalPath)
	tex, err := loadTexture(normalPath, fs)
	if err != nil {
		return nil, err
	}

	textureManager.NewCustom(path, tex)
	textureManager.NewCustom(normalPath, tex)

	return tex, nil
}

func normaliseTexturePath(path string) string {
	// TODO

	return path
}
