package textures

import (
	"fmt"
	"strings"

	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/vtf"
)

type Texture struct {
	filePath   string
	fileSystem *filesystem.Filesystem
	width      int
	height     int
	vtf        *vtf.Vtf
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
