package cameras

import (
	"fmt"
	"math"

	"github.com/emily33901/forgery/fcore"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/math32"
)

type Camera struct {
	*camera.Camera
}

// Y is up

func (c *Camera) Rotate(x, y, z float32) {
	r := c.Camera.Rotation()

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

func (c *Camera) GetRightVector() (*math32.Vector3, *math32.Vector3) {
	rot := c.Camera.Rotation()

	f := math32.NewVector3(
		math32.Cos(rot.Z)*math32.Sin(rot.Z),
		math32.Cos(rot.Z)*math32.Cos(rot.Z),
		math32.Sin(rot.Z))

	r := math32.NewVector3(
		(math32.Sin((rot.Y) - math.Pi/2)),
		(math32.Cos((rot.Y) - math.Pi/2)),
		0,
	)

	return f, r
}

func (c *Camera) Move(forward, back, left, right bool, scale float32) {
	// camForward, _ := c.GetForwardRightVectors()
	// p := c.Position()

	var f, b float32

	if forward {
		fmt.Println("forward")
		f = 1 * scale
	}

	if back {
		fmt.Println("back")
		b = 1 * scale
	}
	// l := int(left)
	// r := int(right)

	d := c.Direction()

	c.Camera.TranslateOnAxis(&d, b-f)

	fmt.Println(c.Camera.Position())
}

var cameras *fcore.Manager = fcore.NewManager("camera-%d")

func New() string {
	c := &Camera{
		camera.New(1),
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
