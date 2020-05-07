package world

import (
	"fmt"

	"github.com/g3n/engine/math32"
)

type Solid struct {
	Id     int
	Sides  []Side
	Editor *Editor
}

type Side struct {
	Id              int
	Plane           Plane
	Material        string
	UAxis           UVTransform
	VAxis           UVTransform
	Rotation        float32
	LightmapScale   float32
	SmoothingGroups bool
}

type UVTransform struct {
	Transform math32.Vector4
	Scale     float32
}

type Editor struct {
	Color             math32.Vector3
	visgroupShown     bool
	visGroupAutoShown bool

	logicalPos math32.Vector2 // only exists on brush entities?
}

type Plane [3]math32.Vector3

func NewSolid(id int, sides []Side, editor *Editor) *Solid {
	return &Solid{
		Id:     id,
		Sides:  sides,
		Editor: editor,
	}
}

func NewSide(id int, plane Plane, material string, uAxis UVTransform, vAxis UVTransform, rotation float32, lightmapScale float32, smoothingGroups bool) *Side {
	return &Side{
		Id:              id,
		Plane:           plane,
		Material:        material,
		UAxis:           uAxis,
		VAxis:           vAxis,
		Rotation:        rotation,
		LightmapScale:   lightmapScale,
		SmoothingGroups: smoothingGroups,
	}
}

func NewEditor(color math32.Vector3, visgroupShown bool, visgroupAutoShown bool) *Editor {
	return &Editor{
		Color:             color,
		visgroupShown:     visgroupShown,
		visGroupAutoShown: visgroupAutoShown,
	}
}

func NewPlane(a math32.Vector3, b math32.Vector3, c math32.Vector3) *Plane {
	p := Plane([3]math32.Vector3{a, b, c})
	return &p
}

func NewPlaneFromString(marshalled string) *Plane {
	var v1, v2, v3 = float32(0), float32(0), float32(0)
	var v4, v5, v6 = float32(0), float32(0), float32(0)
	var v7, v8, v9 = float32(0), float32(0), float32(0)
	fmt.Sscanf(marshalled, "(%f %f %f) (%f %f %f) (%f %f %f)", &v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9)

	return NewPlane(
		math32.Vector3{v1, v2, v3},
		math32.Vector3{v4, v5, v6},
		math32.Vector3{v7, v8, v9})
}

func NewUVTransform(transform math32.Vector4, scale float32) *UVTransform {
	return &UVTransform{
		Transform: transform,
		Scale:     scale,
	}
}

func NewUVTransformFromString(marshalled string) *UVTransform {
	var v1, v2, v3, v4 = float32(0), float32(0), float32(0), float32(0)
	var scale = float32(0)
	fmt.Sscanf(marshalled, "[%f %f %f %f] %f", &v1, &v2, &v3, &v4, &scale)
	return NewUVTransform(math32.Vector4{v1, v2, v3, v4}, scale)
}
