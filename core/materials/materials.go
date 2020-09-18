package materials

import (
	"fmt"

	"github.com/emily33901/forgery/core/events"
	"github.com/emily33901/forgery/core/textures"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/golang-source-engine/vmt"
)

// Material
type Material struct {
	Props *vmt.Properties

	filePath string
	loaded   bool

	Textures struct {
		// Albedo
		albedoPath string
		Albedo     *textures.Texture

		// Normal
		normalPath string
		Normal     *textures.Texture
	}

	g3nMaterial *material.Standard
}

func (mat *Material) Loaded() bool {
	return mat.loaded
}

func (mat *Material) TextureLoaded(ev *textures.TextureLoadedEvent) {
	if ev.Err != nil {
		fmt.Println("Failed to load texture", ev.Err)
	}

	if ev.Path == mat.Textures.albedoPath {
		fmt.Println("Loaded", ev.Path)
		mat.Textures.Albedo = ev.Tex
		mat.loaded = true
	}

	if ev.Path == mat.Textures.normalPath {
		mat.Textures.Normal = ev.Tex
	}
}

var errorMat *material.Standard
var errorTex *texture.Texture2D

func (mat *Material) G3nMaterial() *material.Standard {
	if !mat.loaded {
		if errorMat != nil {
			return errorMat
		}

		errorMat = material.NewStandard(&math32.Color{1, 1, 1})
		errorTex = texture.NewBoard(8, 8, &math32.Color{0, 0, 0}, &math32.Color{1, 0, 1}, &math32.Color{1, 0, 1}, &math32.Color{0, 0, 0}, 1.0)
		errorTex.SetWrapS(gls.REPEAT)
		errorTex.SetWrapT(gls.REPEAT)
		errorTex.SetRepeat(8, 8)
		// errorMat = material.NewBasic().GetMaterial()
		errorMat.AddTexture(errorTex)

		return errorMat
	}

	if mat.g3nMaterial != nil {
		return mat.g3nMaterial
	}

	mat.g3nMaterial = material.NewStandard(&math32.Color{1, 1, 1})
	albedoTex := mat.Textures.Albedo.G3nTexture()
	mat.g3nMaterial.AddTexture(albedoTex)

	if mat.Textures.Albedo.Translucent || (mat.Props.AlphaTest != "" && mat.Props.Alpha < 1.0) {
		mat.g3nMaterial.SetTransparent(mat.Textures.Albedo.Translucent)
		mat.g3nMaterial.SetOpacity(mat.Props.Alpha)
	}

	// TODO when they add AddNormalMap use that instead of this!
	if mat.Textures.Normal != nil {
		normalTex := mat.Textures.Normal.G3nTexture()
		normalTex.SetUniformNames("uNormalSampler", "uNormalTexParams")
		mat.g3nMaterial.ShaderDefines.Set("HAS_NORMALMAP", "")

		mat.g3nMaterial.AddTexture(normalTex)
	}

	return mat.g3nMaterial
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

func NewMaterial(filePath string, props *vmt.Properties) (mat *Material) {
	mat = &Material{
		filePath: filePath,
		Props:    props,
	}

	events.Subscribe(textures.TextureLoaded, func(_ string, evData interface{}) {
		ev := evData.(*textures.TextureLoadedEvent)
		mat.TextureLoaded(ev)
	})

	return
}
