package windows

import (
	"fmt"
	"math"

	"github.com/emily33901/forgery/core/manager"
	"github.com/emily33901/forgery/core/scenes"
	fcore "github.com/emily33901/forgery/forgery/core"
	"github.com/emily33901/forgery/forgery/render"
	"github.com/emily33901/forgery/forgery/render/cameras"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/inkyblackness/imgui-go"
)

type SceneWindow struct {
	core.IDispatcher
	Scene    *scenes.Scene
	cameraId string
	id       string
	closing  bool

	fb          *render.FrameBuffer
	adapter     render.Adapter
	platform    fcore.Platform
	size        imgui.Vec2
	fbPos       imgui.Vec2
	renderSize  imgui.Vec2
	sizeChanged bool

	lastMouseHitPos imgui.Vec2

	focused       bool
	dragging      bool
	lastDragDelta imgui.Vec2
}

var sceneWindows *manager.Manager = manager.NewManager("scenewindow-%d")

func oglToImguiTextureId(id uint32) imgui.TextureID {
	return imgui.TextureID(uint64(id) | (1 << 32))
}

func Iter(cb func(k string, v *SceneWindow)) {
	sceneWindows.Iter(func(k string, v interface{}) {
		cb(k, v.(*SceneWindow))
	})
}

func NewSceneWindow(cameraId string, adapter render.Adapter, platform fcore.Platform) *SceneWindow {
	if cameraId == "" {
		cameraId = cameras.New()
	}

	w := &SceneWindow{
		Scene:    scenes.New(),
		cameraId: cameraId,
		fb:       render.NewFramebuffer(adapter, 200, 200),
		adapter:  adapter,
		platform: platform,
	}

	w.Scene.Root.Add(cameras.Get(cameraId))

	w.id = sceneWindows.New(w)

	w.AttachCameraToScene()

	return w
}

func (w *SceneWindow) bind() {
	w.fb.Bind()
	w.adapter.Viewport(0, 0, int32(w.size.X), int32(w.size.Y))
}

func (w *SceneWindow) unbind() {
	w.fb.Unbind()
}

func (w *SceneWindow) startFrame() {
	w.adapter.ClearAll()
}

func (w *SceneWindow) endFrame() {
}

func (w *SceneWindow) Render(r *renderer.Renderer) {
	if w.sizeChanged {
		w.sizeChanged = false
		w.Camera().SetAspect(w.size.X / w.size.Y)
		w.fb.Resize(int(w.size.X), int(w.size.Y))
	}

	w.bind()
	w.startFrame()
	r.Render(w.Scene.Root, cameras.Get(w.cameraId))
	w.endFrame()
	w.unbind()
}

func (w *SceneWindow) Camera() *cameras.Camera {
	return cameras.Get(w.cameraId)
}

func (w *SceneWindow) AttachCameraToScene() {
	// camera.NewOrbitControl(w.Camera())
}

func (w *SceneWindow) focusedControl(deltaTime float32) {
	pos := imgui.WindowPos()
	windowCentre := pos.Plus(w.size.Times(0.5))

	{
		// Draw a hitcursor
		oldCursorPos := imgui.CursorScreenPos()
		imgui.SetCursorScreenPos(windowCentre)
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{0, 1, 0, 1})
		imgui.Text("+")
		imgui.PopStyleColor()
		imgui.SetCursorScreenPos(oldCursorPos)
	}

	delta := imgui.CurrentIO().MouseDelta().Times(deltaTime)

	//fmt.Println("Delta is", delta)

	w.Camera().Rotate(delta.X, delta.Y, 0)

	w.Camera().Move(
		w.platform.KeyDown('W'), w.platform.KeyDown('S'),
		w.platform.KeyDown('A'), w.platform.KeyDown('D'),
		10*deltaTime)
}

var zeroVec imgui.Vec2 = imgui.Vec2{0, 0}

func (w *SceneWindow) unfocusedControl(deltaTime float32) {
	{
		// Handle dragging

		dragDelta := imgui.MouseDragDelta(0, 10).Times(deltaTime)

		if !w.dragging && dragDelta != zeroVec {
			w.dragging = true
			w.lastDragDelta = zeroVec

			return
		} else if w.dragging && dragDelta == zeroVec {
			w.dragging = false

			return
		} else if w.dragging {
			realDelta := dragDelta.Minus(w.lastDragDelta)
			w.lastDragDelta = dragDelta
			w.Camera().Rotate(realDelta.X, realDelta.Y, 0)

			return
		}
	}

	{
		// Handle selecting an object
		if imgui.IsMouseClicked(0) {
			mouseScreenPos := imgui.MousePos()

			mouseWindowPos := mouseScreenPos.Minus(w.fbPos).Minus(imgui.WindowPos())

			normalisedCoords := imgui.Vec2{
				(-.5 + mouseWindowPos.X/float32(w.size.X)) * 2.0,
				(.5 - mouseWindowPos.Y/float32(w.size.Y)) * 2.0,
			}

			r := collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
			r.SetFromCamera(w.Camera().Camera, normalisedCoords.X, normalisedCoords.Y)

			results := r.IntersectObject(w.Scene.Root, true)

			for _, r := range results {
				fmt.Println("Hit object at", r.Object.Position())

				if r.Object.Position() == (math32.Vector3{0, 0, 0}) {
					r.Object.GetNode().SetPosition(0, 1, 0)
				} else if r.Object.Position() == (math32.Vector3{0, 1, 0}) {
					r.Object.GetNode().SetPosition(0, 0, 0)
				}

				// TODO send object-selected event

			}
			w.lastMouseHitPos = mouseWindowPos
		}

		{
			oldCursorPos := imgui.CursorPos()
			imgui.SetCursorPos(w.lastMouseHitPos.Plus(w.fbPos))
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{1, 0, 0, 1})
			imgui.Text("+")
			imgui.PopStyleColor()

			imgui.SetCursorPos(oldCursorPos)
		}
	}
}

func (w *SceneWindow) BuildUI(deltaTime float32) {
	imgui.SetNextWindowSizeConstraints(imgui.Vec2{100, 100}, imgui.Vec2{math.MaxFloat32, math.MaxFloat32})
	if imgui.BeginV(w.id, &w.closing, imgui.WindowFlagsNoScrollbar) {
		size := imgui.ContentRegionAvail()

		if size != w.size {
			w.size = size
			w.sizeChanged = true
		}

		fbWidth, fbHeight := w.fb.Size()

		w.fbPos = imgui.CursorPos()

		imgui.ImageButtonV(oglToImguiTextureId(w.fb.TextureId()),
			size, //imgui.Vec2{X: wSize.X, Y: wSize.Y},
			imgui.Vec2{X: 0, Y: size.Y / float32(fbHeight)},
			imgui.Vec2{X: size.X / float32(fbWidth), Y: 0},
			0,
			imgui.Vec4{X: 0, Y: 0, Z: 0, W: 1}, imgui.Vec4{X: 1, Y: 1, Z: 1, W: 1})

		if imgui.IsKeyPressedV(int('Z'), false) {
			if !w.focused && imgui.IsItemHovered() {
				fmt.Println("Entered focus")
				w.focused = true
				w.platform.SetCursorEnabled(false)
			} else if w.focused {
				fmt.Println("Exited focus")
				w.focused = false
				w.platform.SetCursorEnabled(true)
			}
		}

		if w.focused {
			w.focusedControl(deltaTime)
		} else if imgui.IsItemHovered() || w.dragging {
			w.unfocusedControl(deltaTime)
		}

	}
	imgui.End()
}
