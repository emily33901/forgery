package cameras

import (
	"github.com/emily33901/forgery-go/core"
	"github.com/g3n/engine/camera"
)

var cameras *core.Manager = core.NewManager("camera-%d")

func New() string {
	c := camera.New(1)

	return cameras.New(c)
}

func Get(c string) *camera.Camera {
	v := cameras.Get(c)

	if v == nil {
		return nil
	}

	return v.(*camera.Camera)
}
