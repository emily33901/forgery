package materials

import (
	"fmt"

	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/forgery/core/manager"
	"github.com/emily33901/forgery/core/textures"
	"github.com/golang-source-engine/vmt"
)

// TODO we might aswell just instantly convert into g3ns material type rather than
// trying to keep any semblance of sources stuff here

var materialsManager = manager.NewManager("")

func Load(path string, fs *filesystem.Filesystem) (*Material, error) {
	v := materialsManager.Get(path)

	if v != nil {
		return v.(*Material), nil
	}

	normalPath := normaliseMaterialPath(path)

	// Try the normalised path aswell
	v = materialsManager.Get(normalPath)

	if v != nil {
		return v.(*Material), nil
	}

	// attempt to load
	fmt.Println("Loading", normalPath)

	mat, err := loadMaterial(normalPath, fs)

	if err != nil {
		return nil, err
	}

	materialsManager.NewCustom(path, mat)
	materialsManager.NewCustom(normalPath, mat)

	return mat, nil
}

func normaliseMaterialPath(path string) string {
	// TODO

	return path
}

func loadMaterial(path string, fs *filesystem.Filesystem) (*Material, error) {
	vtfTexturePath := ""

	vmtProperties, err := vmt.FromFilesystem(path, fs, vmt.NewProperties())

	if err != nil {
		fmt.Println("Failed to load material", path, "Reason", err)
		return nil, err
	}

	mat := NewMaterial(path, vmtProperties.(*vmt.Properties))
	if mat.Props.BaseTexture == "" {
		// Has no texture so load error texture
		// TODO
	}

	vtfTexturePath = mat.Props.BaseTexture
	mat.Textures.Albedo, err = textures.Load(vtfTexturePath, fs)

	if err != nil {
		// Didnt find a texture so make it the error texture
		fmt.Println("Failed to load texture", vtfTexturePath, err)
		// TODO
	}

	if mat.Props.Bumpmap != "" {
		mat.Textures.Normal, err = textures.Load(mat.Props.Bumpmap, fs)

		if err != nil {
			fmt.Println("Failed to load normal texture even though there is one")
		}

		// TODO error texture?
	}

	return mat, nil
}
