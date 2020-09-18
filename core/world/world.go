package world

import (
	"fmt"
	"strconv"

	"github.com/emily33901/forgery/core/events"
	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/forgery/core/materials"
	"github.com/emily33901/forgery/core/textures"
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

	events.Subscribe(textures.TextureLoaded, func(_ string, evdata interface{}) {
		w.MakeDirty()
	})

	return w
}

func (w *World) MakeDirty() {
	w.sceneDirty = true
}

func CreateFace(w *Winding, materialName string, fs *filesystem.Filesystem, uaxis, vaxis *UVTransform) (*geometry.Geometry, *material.Standard) {
	geom := geometry.NewGeometry()

	// Get the verticies
	verts := math32.NewArrayF32(0, 16)

	flipVertex := func(v ...*math32.Vector3) (ret []*math32.Vector3) {
		for _, x := range v {
			ret = append(ret, &math32.Vector3{
				x.X,
				x.Z,
				x.Y,
			})
		}
		return
	}

	// flipTransform := func(v ...*UVTransform) (ret []*UVTransform) {
	// 	for _, x := range v {
	// 		ret = append(ret, &UVTransform{
	// 			Transform: math32.Vector4{
	// 				x.Transform.X,
	// 				x.Transform.Z,
	// 				x.Transform.Y,
	// 				x.Transform.W,
	// 			},
	// 			Scale: x.Scale,
	// 		})
	// 	}
	// 	return
	// }

	// flipVertex := func(v ...*math32.Vector3) (ret []*math32.Vector3) {
	// 	ret = v
	// 	return
	// }

	verts.AppendVector3(
		flipVertex(w.Points...)...,
	)

	// Create a vbo
	indicies := math32.NewArrayU32(0, verts.Len()*2)
	i := 0
	j := 1
	for j+1 < len(w.Points) {
		indicies.Append(0)
		indicies.Append(uint32(j))
		indicies.Append(uint32(j + 1))

		j += 1
		i += 3
	}

	// normals
	normals := math32.NewArrayF32(0, 16)

	// Calculate the normal by picking 2 points and crossing them
	a := flipVertex(w.Points[0])[0]
	b := flipVertex(w.Points[1])[0].Sub(a)
	c := flipVertex(w.Points[2])[0].Sub(a)
	normalVec := b.Cross(c)

	for i := 0; i < indicies.Len(); i++ {
		normals.AppendVector3(normalVec)
	}

	// uvs
	sourceMat, err := materials.Load(materialName, fs)

	if err != nil {
		panic(err)
	}

	// mat := material.NewStandard(&math32.Color{1, 1, 1})
	mat := sourceMat.G3nMaterial()
	mat.SetUseLights(material.UseLightAll)
	// mat.SetTransparent(false)
	// mat.SetWireframe(false)
	// all sides face inwards so we want to draw the back face
	mat.SetSide(material.SideBack)

	width := 128
	height := 128

	if sourceMat.Loaded() {
		width = sourceMat.Width()
		height = sourceMat.Height()
	}

	u := uaxis
	v := vaxis

	uvs := math32.NewArrayF32(0, 16)
	for _, vertex := range w.Points {
		cu := ((u.Transform.X * vertex.X) +
			(u.Transform.Y * vertex.Y) +
			(u.Transform.Z * vertex.Z)) / float32(u.Scale) / float32(width)

		cv := ((v.Transform.X * vertex.X) +
			(v.Transform.Y * vertex.Y) +
			(v.Transform.Z * vertex.Z)) / float32(v.Scale) / float32(height)

		// uvs.Append(1.0, 1.0)
		uvs.Append(cu, cv)
	}

	geom.SetIndices(indicies)
	geom.AddVBO(gls.NewVBO(verts).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	geom.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))
	// gls.NewVBO(verts).AddAttrib(gls.VertexTexcoord)

	return geom, mat
}

// BuildScene converts the internal representation into
// a g3n scene which can be rendered
func (w *World) BuildScene(fs *filesystem.Filesystem) {
	if w.sceneDirty == false {
		return
	}

	// Cleanup the old scene
	w.SceneSolid.DisposeChildren(true)
	w.SceneSolid.SetName("World main node")
	w.SceneWireframe.DisposeChildren(true)
	w.Root.Add(helper.NewAxes(128))

	// Start building a new scene
	// This essentially goes through every solid and whatnot
	// building up nodes out of geometry

	for _, s := range w.solids {
		solidNode := core.NewNode()
		w.SceneSolid.Add(solidNode)

		solidNode.SetLoaderID(strconv.Itoa(s.Id))

		// https://github.com/emily33901/HammerFromScratch/blob/a0f669718a70632138545fd1a5a493b8299221a0/hammer/mapsolid.cpp#L788

		usePlane := make([]bool, len(s.Sides))

		for i, side := range s.Sides {
			if side.Plane.Normal.LengthSq() == 0 {
				// Not a valid plane
				usePlane[i] = false
				continue
			}

			usePlane[i] = true

			// Check this plane isnt identical to another plane
			for j, side2 := range s.Sides {
				if i == j {
					break
				}

				if side.Plane.Normal.Dot(&side2.Plane.Normal) > 0.999 && math32.Abs(side.Plane.Dist-side2.Plane.Dist) < 0.1 {
					usePlane[j] = false
				}
			}
		}

		// Now that we have all of the faces and we know which ones to use
		// its time to clip all of them to get the points

		for i, side := range s.Sides {
			if usePlane[i] == false {
				// we are not using this plane
				continue
			}

			winding := CreateWindingFromPlane(&side.Plane)

			for j, side := range s.Sides {
				if j != i && len(winding.Points) > 0 {
					winding.Clip(&side.Plane)
				}
			}

			if len(winding.Points) == 0 {
				fmt.Println("Empty winding")
				continue
			}

			// This winding has points in it that can be turned into a face
			geom, mat := CreateFace(winding, side.Material, fs, &side.UAxis, &side.VAxis)

			sideNode := graphic.NewMesh(geom, mat)
			sideNode.SetVisible(true)
			sideNode.SetLoaderID(strconv.Itoa(side.Id))

			solidNode.Add(sideNode)
		}

		// TODO wireframe
	}

	l1 := light.NewAmbient(&math32.Color{1, 1, 1}, 1.0)
	w.Root.Add(l1)

	w.sceneDirty = false
}
