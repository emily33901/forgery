package keyvalues

import (
	"os"

	keyvalues "github.com/galaco/KeyValues"
)

func FromDisk(path string) (*keyvalues.KeyValue, error) {
	stream, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	kvReader := keyvalues.NewReader(stream)

	kv, err := kvReader.Read()

	return &kv, err
}
