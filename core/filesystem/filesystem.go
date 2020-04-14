package filesystem

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/galaco/bsp/lumps"
	vpk "github.com/galaco/vpk2"
)

type Filesystem struct {
	gameVPKs         map[string]vpk.VPK
	localDirectories []string
	pakFile          *lumps.Pakfile
}

func NewFilesystem() *Filesystem {
	return &Filesystem{
		gameVPKs:         map[string]vpk.VPK{},
		localDirectories: []string{},
		pakFile:          nil,
	}
}

// PakFile returns loaded pakfile
// There can only be 1 registered pakfile at once.
func (fs *Filesystem) PakFile() *lumps.Pakfile {
	return fs.pakFile
}

// RegisterVpk registers a vpk package as a valid
// asset directory
func (fs *Filesystem) RegisterVpk(path string, vpkFile *vpk.VPK) {
	fs.gameVPKs[path] = *vpkFile
}

func (fs *Filesystem) UnregisterVpk(path string) {
	for key := range fs.gameVPKs {
		if key == path {
			delete(fs.gameVPKs, key)
		}
	}
}

// RegisterLocalDirectory register a filesystem path as a valid
// asset directory
func (fs *Filesystem) RegisterLocalDirectory(directory string) {
	fs.localDirectories = append(fs.localDirectories, directory)
}

func (fs *Filesystem) UnregisterLocalDirectory(directory string) {
	for idx, dir := range fs.localDirectories {
		if dir == directory {
			if len(fs.localDirectories) == 1 {
				fs.localDirectories = make([]string, 0)
				return
			}
			fs.localDirectories = append(fs.localDirectories[:idx], fs.localDirectories[idx+1:]...)
		}
	}
}

// RegisterPakFile Set a pakfile to be used as an asset directory.
// This would normally be called during each map load
func (fs *Filesystem) RegisterPakFile(pakFile *lumps.Pakfile) {
	fs.pakFile = pakFile
}

// UnregisterPakFile removes the current pakfile from
// available search locations
func (fs *Filesystem) UnregisterPakFile() {
	fs.pakFile = nil
}

// EnumerateResourcePaths returns all registered resource paths.
// PakFile is excluded.
func (fs *Filesystem) EnumerateResourcePaths() []string {
	list := make([]string, 0)

	for idx := range fs.gameVPKs {
		list = append(list, string(idx))
	}

	list = append(list, fs.localDirectories...)

	return list
}

// GetFile attempts to get stream for filename.
// Search order is Pak->FileSystem->VPK
func (fs *Filesystem) GetFile(filename string) (io.Reader, error) {
	// sanitise file
	searchPath := NormalisePath(strings.ToLower(filename))

	// try to read from pakfile first
	if fs.pakFile != nil {
		f, err := fs.pakFile.GetFile(searchPath)
		if err == nil && f != nil && len(f) != 0 {
			return bytes.NewReader(f), nil
		}
	}

	// Fallback to local filesystem
	for _, dir := range fs.localDirectories {
		if _, err := os.Stat(dir + "\\" + searchPath); os.IsNotExist(err) {
			continue
		}
		file, err := ioutil.ReadFile(dir + searchPath)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(file), nil
	}

	// Fall back to game vpk
	for _, vfs := range fs.gameVPKs {
		entry := vfs.Entry(searchPath)
		if entry != nil {
			return entry.Open()
		}
	}

	return nil, NewFileNotFoundError(filename)
}

// AllPaths returns all the paths (files) that are currently loaded
func (fs *Filesystem) AllPaths() []string {
	results := []string{}

	// Start with the local directories
	for _, dir := range fs.localDirectories {
		finfo, err := ioutil.ReadDir(dir)
		if err == nil {
			for _, f := range finfo {
				// TODO check this returns the correct file name
				results = append(results, f.Name())
			}
		}
	}

	// Now do vpks
	for _, vfs := range fs.gameVPKs {
		results = append(results, vfs.Paths()...)
	}

	// TODO: handle pak files

	return results
}
