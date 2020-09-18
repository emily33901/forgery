package textures

import (
	"fmt"

	"github.com/emily33901/forgery/core/events"
	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/forgery/core/manager"
)

const (
	// Events

	// TextureLoaded tells other systems that a texture is loaded
	TextureLoaded = "Textures.Loaded"
)

type TextureLoadedEvent struct {
	Path string
	Tex  *Texture
	Err  error
}

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
		events.Dispatch(TextureLoaded, &TextureLoadedEvent{
			path, nil, err,
		})

		return nil, err
	}

	textureManager.NewCustom(path, tex)
	textureManager.NewCustom(normalPath, tex)

	// Tell other people that we loaded it
	events.Dispatch(TextureLoaded, &TextureLoadedEvent{
		path, tex, nil,
	})
	return tex, nil
}

func normaliseTexturePath(path string) string {
	// TODO

	return path
}
