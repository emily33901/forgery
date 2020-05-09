package cameras

import (
	"github.com/emily33901/forgery/core/manager"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/math32"
)

type Camera struct {
	*camera.Camera
}

// Y is up

func (c *Camera) Rotate(x, y, z float32) {
	c.RotateFrom(c.Rotation(), x, y, z)

}

func (c *Camera) RotateFrom(r math32.Vector3, x, y, z float32) {
	r.Y += x
	r.X -= y

	if r.X < math32.DegToRad(90) {
		r.X = math32.DegToRad(90)
	}

	if r.X > math32.DegToRad(270) {
		r.X = math32.DegToRad(270)
	}

	c.Camera.SetRotation(r.X, r.Y, math32.DegToRad(180))
}

func (c *Camera) GetForwardRightVector() (forward math32.Vector3, right *math32.Vector3) {
	forward = c.Camera.Direction()
	up := &math32.Vector3{0, 1, 0}

	right = forward.Clone().Cross(up)

	return
}

func (c *Camera) Move(forward, back, left, right bool, scale float32) {
	var f, b, l, r float32

	if forward {
		f = 1 * scale
	}

	if back {
		b = 1 * scale
	}

	if left {
		l = 1 * scale
	}

	if right {
		r = 1 * scale
	}

	d, dr := c.GetForwardRightVector()

	c.Camera.TranslateOnAxis(&d, b-f)
	c.Camera.TranslateOnAxis(dr, l-r)
}

var cameras *manager.Manager = manager.NewManager("camera-%d")

func New() string {
	c := &Camera{
		camera.NewPerspective(1, 0.3, 8192, 90, camera.Horizontal),
	}

	c.Camera.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 1, 0})

	return cameras.New(c)
}

func Get(c string) *Camera {
	v := cameras.Get(c)

	if v == nil {
		return nil
	}

	return v.(*Camera)
}

func Iter(cb func(k string, v *Camera)) {
	cameras.Iter(func(k string, v interface{}) {
		cb(k, v.(*Camera))
	})
}

func Keys() []string {
	return cameras.Keys()
}
