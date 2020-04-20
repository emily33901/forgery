package scenes

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Scene maps a map file (vmf) and a g3n scene together
// (it technically doesnt do this but it keeps track of
//  geomentry and similar and allows for exporting back
//  to a vmf)

type Scene struct {
	Geom  *core.Node
	Root  *core.Node
	Debug *core.Node
	// Node *core.Node
}

func New() *Scene {
	s := &Scene{}

	s.Root = core.NewNode()
	s.Geom = s.Root.Add(core.NewNode())
	s.Debug = s.Root.Add(core.NewNode())

	geom := geometry.NewPlane(10, 10)

	mesh := graphic.NewMesh(geom, material.NewStandard(math32.NewColorHex(0x55331122)))
	mesh.SetPosition(0, 0, 0)
	mesh.SetRotation(10, 10, 10)

	s.Geom.Add(mesh)

	return s
}
