package textures

import (
	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/forgery/core/manager"
)

var textureManager *manager.Manager = manager.NewManager("")

func Load(path string, fs *filesystem.Filesystem) (*Texture, error) {
	v := textureManager.Get(path)

	if v == nil {
		// attempt to load
		tex, err := loadTexture(path, fs)
		if err != nil {
			return nil, err
		}

		return tex, nil
	}

	return v.(*Texture), nil
}
