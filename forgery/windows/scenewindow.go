package windows

import (
	"fmt"
	"math"

	"github.com/emily33901/forgery/core/manager"
	fcore "github.com/emily33901/forgery/forgery/core"
	"github.com/emily33901/forgery/forgery/render"
	"github.com/emily33901/forgery/forgery/render/cameras"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/renderer"
	"github.com/inkyblackness/imgui-go"
)

type SceneWindow struct {
	core.IDispatcher
	Scene    *core.Node
	cameraId string
	id       string
	closing  bool

	fb          *render.FrameBuffer
	adapter     render.Adapter
	platform    fcore.Platform
	size        imgui.Vec2
	sizeChanged bool

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
		Scene:    core.NewNode(),
		cameraId: cameraId,
		fb:       render.NewFramebuffer(adapter, 200, 200),
		adapter:  adapter,
		platform: platform,
	}

	w.Scene.Add(cameras.Get(cameraId))

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
	r.Render(w.Scene, cameras.Get(w.cameraId))
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
	dragDelta := imgui.MouseDragDelta(0, 10).Times(deltaTime)

	if !w.dragging && dragDelta != zeroVec {
		w.dragging = true
		w.lastDragDelta = zeroVec
	} else if w.dragging && dragDelta == zeroVec {
		w.dragging = false
	} else if w.dragging {
		realDelta := dragDelta.Minus(w.lastDragDelta)
		w.lastDragDelta = dragDelta
		w.Camera().Rotate(realDelta.X, realDelta.Y, 0)
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

		fbWidth, hbHeight := w.fb.Size()

		imgui.ImageButtonV(oglToImguiTextureId(w.fb.TextureId()),
			size, //imgui.Vec2{X: wSize.X, Y: wSize.Y},
			imgui.Vec2{X: 0, Y: size.Y / float32(hbHeight)},
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