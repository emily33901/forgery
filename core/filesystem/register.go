package filesystem

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	keyvalues "github.com/galaco/KeyValues"
	vpk "github.com/galaco/vpk2"
)

// CreateFilesystemFromGameInfoDefinitions Reads game resource data paths
// from gameinfo.txt
// All games should ship with a gameinfo.txt, but it isn't actually mandatory.
func CreateFromGameInfo(basePath string, gameInfo *keyvalues.KeyValue) *Filesystem {
	fs := NewFilesystem()
	fsNode, _ := gameInfo.Find("FileSystem")

	searchPathsNode, _ := fsNode.Find("SearchPaths")
	searchPaths, _ := searchPathsNode.Children()
	basePath, _ = filepath.Abs(basePath)
	basePath = strings.Replace(basePath, "\\", "/", -1)

	for _, searchPath := range searchPaths {
		kv := searchPath
		path, _ := kv.AsString()
		path = strings.Trim(path, " ")

		// Current directory
		gameInfoPathRegex := regexp.MustCompile(`(?i)\|gameinfo_path\|`)
		if gameInfoPathRegex.MatchString(path) {
			path = gameInfoPathRegex.ReplaceAllString(path, basePath+"/")
		}

		// Executable directory
		allSourceEnginePathsRegex := regexp.MustCompile(`(?i)\|all_source_engine_paths\|`)
		if allSourceEnginePathsRegex.MatchString(path) {
			path = allSourceEnginePathsRegex.ReplaceAllString(path, basePath+"/../")
		}
		if strings.Contains(strings.ToLower(kv.Key()), "mod") && !strings.HasPrefix(path, basePath) {
			path = basePath + "/../" + path
		}

		// Strip vpk extension, then load it
		path = strings.Trim(strings.Trim(path, " "), "\"")
		if strings.HasSuffix(path, ".vpk") {
			path = strings.Replace(path, ".vpk", "", 1)
			opener := vpk.MultiVPK(path)
			vpkHandle, err := vpk.Open(opener)
			if err != nil {
				// TODO log error
				continue
			}
			fs.RegisterVpk(path, vpkHandle)
		} else {
			// wildcard suffixes not useful
			if strings.HasSuffix(path, "/*") {
				path = strings.Replace(path, "/*", "", -1)
			}
			fs.RegisterLocalDirectory(path)
		}
	}

	return fs
}

func CreateFromGameDir(path string, gameInfo *keyvalues.KeyValue) *Filesystem {

	// Register GameInfo.txt referenced resource paths
	// Filesystem module needs to know about all the possible resource
	// locations it can search.
	fs := CreateFromGameInfo(path, gameInfo)

	// Make sure to also load the platform dir
	fs.RegisterLocalDirectory(path + "/../platform")

	// Now try and load all of the vpks that are in those directories
	for _, x := range fs.EnumerateResourcePaths() {
		files, err := ioutil.ReadDir(x)

		if err != nil {
			// panic(err)
			continue
		}

		for _, f := range files {
			if strings.HasSuffix(f.Name(), "_dir.vpk") {
				nameNoSuffix := x + "/" + f.Name()[:len(f.Name())-8]
				opener := vpk.MultiVPK(nameNoSuffix)
				v, err := vpk.Open(opener)

				if err != nil {
					panic(err)
				}

				fs.RegisterVpk(nameNoSuffix, v)
			}
		}

		// TODO: we also need to load all the vpks that arent part of a dir pack
	}

	return fs
}
