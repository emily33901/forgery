package textures

import (
	"fmt"

	"github.com/emily33901/forgery/core/events"
	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/forgery/core/manager"
)

const (
	// Events

	// LoadTexture tells the texture manager to load a texture by name
	LoadTexture = "Textures.LoadTexture"

	// TextureLoaded tells other systems that a texture is loaded
	TextureLoaded = "Textures.Loaded"
)

type LoadTextureEvent struct {
	name string
	fs   *filesystem.Filesystem
}

type TextureLoadedEvent struct {
	name string
	tex  *Texture
}

var textureManager *manager.Manager = manager.NewManager("")

func Init() {
	events.Subscribe(LoadTexture, func(_ string, evData interface{}) {
		ev := evData.(*LoadTextureEvent)
		Load(ev.name, ev.fs)
	})
}

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
