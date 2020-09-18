package textures

import (
	"fmt"
	"strings"

	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/vtf"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/go-gl/gl/v2.1/gl"
)

type Texture struct {
	filePath   string
	fileSystem *filesystem.Filesystem
	width      int
	height     int
	vtf        *vtf.Vtf

	Translucent bool
	Bump        math32.Vector3
	BumpScale   float32
	g3nTexture  *texture.Texture2D
}

// FilePath Get the filepath this data was loaded from
func (tex *Texture) FilePath() string {
	return tex.filePath
}

// Width returns materials width
func (tex *Texture) Width() int {
	return tex.width
}

// Height returns materials height
func (tex *Texture) Height() int {
	return tex.height
}

// Format returns this materials colour format
func (tex *Texture) Format() uint32 {
	if tex.vtf == nil {
		panic("Didnt Reload() before attempting to use a texture")
	}

	return tex.vtf.Header().HighResImageFormat
}

// PixelDataForFrame get raw colour data for this frame
func (tex *Texture) PixelDataForFrame(frame int) []byte {
	if tex.vtf == nil {
		panic("Always Reload() a texture before attempting to access its fields")
	}

	return tex.vtf.HighestResolutionImageForFrame(frame)
}

// Thumbnail returns a small thumbnail image of a material
func (tex *Texture) Thumbnail() []byte {
	if tex.vtf == nil {
		panic("Always Reload() a texture before attempting to access its fields")
	}

	return tex.vtf.LowResImageData()
}

func (tex *Texture) Reload() error {
	stream, err := tex.fileSystem.GetFile(tex.filePath)
	if err != nil {
		fmt.Printf("Unable to load %s from Disk: %s\n", tex.filePath, err)
		return err
	}

	// Attempt to parse the vtf into color data we can use,
	// if this fails (it shouldn't) we can treat it like it was missing
	read, err := vtf.ReadFromStream(stream)
	if err != nil {
		fmt.Printf("Unable to load %s from Disk: %s\n", tex.filePath, err)
		return err
	}

	tex.vtf = read
	return nil
}

func (tex *Texture) G3nTexture() *texture.Texture2D {
	if tex.g3nTexture != nil {
		return tex.g3nTexture
	}

	tex.Reload()
	if isPixelFormatCompressed(glTextureFormatFromVtfFormat(tex.Format())) {
		pxData := tex.PixelDataForFrame(0)
		tex.g3nTexture = texture.NewTexture2DFromCompressedData(tex.width, tex.height, int32(glTextureFormatFromVtfFormat(tex.Format())), int32(len(pxData)), pxData)
	} else {
		tex.g3nTexture = texture.NewTexture2DFromData(tex.width, tex.height, gls.RGBA, gl.UNSIGNED_BYTE, glTextureFormatFromVtfFormat(tex.Format()), tex.PixelDataForFrame(0))
	}

	// Check wether the texture uses alpha
	if tex.vtf.Header().Flags&0x2000 == 0x2000 || tex.vtf.Header().Flags&0x1000 == 0x1000 {
		tex.Translucent = true
	}

	tex.Bump = math32.Vector3{
		tex.vtf.Header().Reflectivity[0],
		tex.vtf.Header().Reflectivity[1],
		tex.vtf.Header().Reflectivity[2],
	}

	tex.BumpScale = tex.vtf.Header().BumpmapScale

	tex.g3nTexture.SetWrapS(gls.REPEAT)
	tex.g3nTexture.SetWrapT(gls.REPEAT)
	tex.g3nTexture.SetRepeat(1, 1)
	tex.g3nTexture.SetFlipY(true)

	tex.EvictFromMainMemory()

	return tex.g3nTexture
}

func (tex *Texture) EvictFromMainMemory() {
	// This will trigger gc to evict this memory
	tex.vtf = nil
}

// NewTexture2D returns a new texture from Vtf
func newTexture(filePath string, fs *filesystem.Filesystem, width int, height int) *Texture {
	// TODO: we should be able to load the vtf all by ourselves!
	return &Texture{
		fileSystem: fs,
		filePath:   filePath,
		width:      width,
		height:     height,
	}
}

func loadTexture(filePath string, fs *filesystem.Filesystem) (*Texture, error) {
	filePath = filesystem.NormalisePath(filePath)
	if !strings.HasSuffix(filePath, filesystem.ExtensionVtf) {
		filePath = filePath + filesystem.ExtensionVtf
	}

	mat, err := readVtf(filesystem.BasePathMaterial+filePath, fs)
	return mat, err
}

// readVtf
func readVtf(path string, fs *filesystem.Filesystem) (*Texture, error) {
	stream, err := fs.GetFile(path)
	if err != nil {
		return nil, err
	}

	header, err := vtf.ReadHeaderFromStream(stream)
	if err != nil {
		return nil, err
	}

	return newTexture(path, fs, int(header.Width), int(header.Height)), nil
}

// gLTextureFormat swap vtf format to openGL format
func glTextureFormatFromVtfFormat(vtfFormat uint32) int {
	switch vtfFormat {
	case 0:
		return gl.RGBA
	case 2:
		return gl.RGB
	case 3:
		return gl.BGR
	case 12:
		return gl.BGRA
	case 13:
		return gl.COMPRESSED_RGB_S3TC_DXT1_EXT
	case 14:
		return gl.COMPRESSED_RGBA_S3TC_DXT3_EXT
	case 15:
		return gl.COMPRESSED_RGBA_S3TC_DXT5_EXT
	default:
		return gl.RGB
	}
}

func isPixelFormatCompressed(format int) bool {
	switch format {
	// , DXT1A
	case gl.COMPRESSED_RGB_S3TC_DXT1_EXT, gl.COMPRESSED_RGBA_S3TC_DXT3_EXT, gl.COMPRESSED_RGBA_S3TC_DXT5_EXT:
		return true
	}

	return false
}
