package world

import (
	"math/rand"
	"strconv"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/helper"
	"github.com/galaco/vmf"
)

type World struct {
	Root *core.Node

	// Always keep a copy of solid and wireframe
	// because we need to raytrace against the solid one
	SceneSolid     *core.Node
	SceneWireframe *core.Node
	// Debug *core.Node

	// Original vmf file if it exists
	vmfFile *vmf.Vmf

	solids     []Solid
	sceneDirty bool
}

func New(solids []Solid) *World {
	w := &World{}

	w.Root = core.NewNode()
	w.solids = solids
	w.SceneSolid = core.NewNode()
	w.SceneWireframe = core.NewNode()
	w.sceneDirty = true

	w.Root.Add(w.SceneSolid).Add(w.SceneWireframe)

	return w
}

// BuildScene converts the internal representation into
// a g3n scene which can be rendered
func (w *World) BuildScene() {
	if w.sceneDirty == false {
		return
	}

	// Cleanup the old scene
	w.SceneSolid.DisposeChildren(true)
	w.SceneSolid.SetName("World core")
	w.SceneWireframe.DisposeChildren(true)
	w.Root.Add(helper.NewAxes(1))

	// Start building a new scene
	// This essentially goes through every solid and whatnot
	// building up nodes out of geometry

	for _, s := range w.solids {
		randomColor := &math32.Color{math32.Mod(rand.Float32(), 0.5) + 0.5, math32.Mod(rand.Float32(), 0.5) + 0.5, math32.Mod(rand.Float32(), 0.5) + 0.5}

		solidNode := core.NewNode()
		w.SceneSolid.Add(solidNode)

		solidNode.SetLoaderID(strconv.Itoa(s.Id))

		for _, side := range s.Sides {
			// Get the verticies
			verts := math32.NewArrayF32(0, 18)
			// a plane represents 3 vertices- bottom-left, top-left and top-right
			// Triangle 1

			flipVertex := func(v *math32.Vector3) (ret *math32.Vector3) {
				ret = v.Clone()
				ret.Y = v.Z
				ret.Z = v.Y

				return
			}

			verts.AppendVector3(
				flipVertex(&side.Plane[0]),
				flipVertex(&side.Plane[1]),
				flipVertex(&side.Plane[2]),

				flipVertex(&side.Plane[0]),
				flipVertex(&side.Plane[2]),

				flipVertex(side.Plane[2].Clone().Sub(side.Plane[1].Clone().Sub(&side.Plane[0]))),
			)

			geom := geometry.NewGeometry()
			geom.AddVBO(gls.NewVBO(verts).AddAttrib(gls.VertexPosition))

			// TODO use actual material here
			// mat := material.NewStandard(&math32.Color{s.Editor.Color.X, s.Editor.Color.Y, s.Editor.Color.Z})
			mat := material.NewStandard(randomColor)
			mat.SetUseLights(material.UseLightAll)

			sideNode := graphic.NewMesh(geom, mat)
			sideNode.SetLoaderID(strconv.Itoa(side.Id))

			solidNode.Add(sideNode)

			// TODO wireframe
		}
	}

	l1 := light.NewAmbient(&math32.Color{1, 1, 1}, 0.5)

	l1.SetPosition(0, 0, 100)
	w.Root.Add(l1)

	w.sceneDirty = false
}
