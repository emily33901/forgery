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
//  geometry and similar and allows for exporting back
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

	// geom := geometry.NewBox(10, 10, 10)

	// box := geom.BoundingBox()

	// var boxSize math32.Vector3
	// box.Size(&boxSize)

	// lineMaterial := material.NewStandard(math32.NewColorHex(0xAABBCC))

	// lineMaterial.SetSide(material.SideDouble)
	// // lineMaterial.SetUseLights(material.)
	// // lineMaterial.SetWireframe(true)

	// // Get all the edges of the box
	// pointList := []*math32.Vector3{
	// 	// 0
	// 	&box.Min,

	// 	// 1
	// 	box.Min.Clone().Add(&math32.Vector3{boxSize.X, 0, 0}),
	// 	box.Min.Clone().Add(&math32.Vector3{0, boxSize.Y, 0}),
	// 	box.Min.Clone().Add(&math32.Vector3{0, 0, boxSize.Z}),

	// 	//4
	// 	box.Min.Clone().Add(&math32.Vector3{boxSize.X, boxSize.Y, 0}),
	// 	box.Min.Clone().Add(&math32.Vector3{0, boxSize.Y, boxSize.Z}),
	// 	box.Min.Clone().Add(&math32.Vector3{boxSize.X, 0, boxSize.Z}),

	// 	// 7
	// 	&box.Max,
	// }

	// edgeList := [][]int{
	// 	{0, 1},
	// 	{0, 2},
	// 	{0, 3},

	// 	{1, 4},
	// 	{1, 6},

	// 	{2, 4},
	// 	{2, 5},

	// 	{3, 5},
	// 	{3, 6},

	// 	{4, 7},
	// 	{5, 7},
	// 	{6, 7},
	// }

	// // Now create the lines
	// lines := []*graphic.Lines{}
	// for _, e := range edgeList {
	// 	p1 := pointList[e[0]]
	// 	p2 := pointList[e[1]]

	// 	lineGeom := NewLine(p1.DistanceTo(p2))

	// 	l := graphic.NewLines(lineGeom, material.NewBasic())

	// 	l.SetPositionVec(p1.Clone().Add(&math32.Vector3{-boxSize.X, -boxSize.Y, -boxSize.Z}))
	// 	l.SetRotationVec(p1.Clone().Sub(p2).Normalize())

	// 	lines = append(lines, l)
	// }

	// for _, l := range lines {
	// 	s.Geom.Add(l)
	// }

	p := geometry.NewPlane(0, 10)

	defaultMat := material.NewStandard(math32.NewColorHex(0xAABBCC))
	defaultMat.SetWireframe(true)
	defaultMat.SetSide(material.SideDouble)

	graphic.NewLines(p, material.NewStandard(math32.NewColorHex(0xAABBCC)))

	mesh := graphic.NewMesh(p, defaultMat)
	mesh.SetPosition(0, 0, 0)
	mesh.SetRotation(10, 10, 10)

	s.Geom.Add(mesh)

	return s
}

func NewLine(length float32) (line *geometry.Geometry) {
	if length <= 0 {
		panic("Length <= 0")
	}

	line = geometry.NewPlane(0, length)
	return
}
