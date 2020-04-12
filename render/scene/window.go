package scene

import (
	"math"

	manager "github.com/emily33901/forgery/core"
	"github.com/emily33901/forgery/render"
	"github.com/emily33901/forgery/render/cameras"
	"github.com/inkyblackness/imgui-go"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/renderer"
)

type Window struct {
	core.IDispatcher
	Scene    *core.Node
	cameraId string
	id       string
	closing  bool

	fb       *render.FrameBuffer
	adapter  render.Adapter
	lastSize imgui.Vec2
}

var windows *manager.Manager = manager.NewManager("window-%d")

func oglToImguiTextureId(id uint32) imgui.TextureID {
	return imgui.TextureID(uint64(id) | (1 << 32))
}

func Iter(cb func(k string, v *Window)) {
	windows.Iter(func(k string, v interface{}) {
		cb(k, v.(*Window))
	})
}

func NewWindow(cameraId string, adapter render.Adapter) *Window {
	if cameraId == "" {
		cameraId = cameras.New()
	}

	w := &Window{
		Scene:    core.NewNode(),
		cameraId: cameraId,
		fb:       render.NewFramebuffer(adapter, 4000, 4000),
		adapter:  adapter,
	}

	w.Scene.Add(cameras.Get(cameraId))

	w.id = windows.New(w)

	w.AttachCameraToScene()

	return w
}

func (w *Window) bind() {
	w.fb.Bind()
	w.adapter.Viewport(0, 0, int32(w.lastSize.X), int32(w.lastSize.Y))
}

func (w *Window) unbind() {
	w.fb.Unbind()
}

func (w *Window) startFrame() {
	w.adapter.ClearAll()
}

func (w *Window) endFrame() {
}

func (w *Window) Render(r *renderer.Renderer) {
	w.bind()
	w.startFrame()
	r.Render(w.Scene, cameras.Get(w.cameraId))
	w.endFrame()
	w.unbind()
}

func (w *Window) Camera() *camera.Camera {
	return cameras.Get(w.cameraId)
}

func (w *Window) AttachCameraToScene() {
	gui.Manager().Set(w.Scene)
	camera.NewOrbitControl(w.Camera())
}

func (w *Window) BuildUI() {
	imgui.SetNextWindowSizeConstraints(imgui.Vec2{100, 100}, imgui.Vec2{math.MaxFloat32, math.MaxFloat32})
	if imgui.BeginV(w.id, &w.closing, 0) {
		size := imgui.ContentRegionAvail()

		if size != w.lastSize {
			w.lastSize = size
			w.Camera().SetAspect(size.X / size.Y)
			// w.fb.Resize(int(size.X), int(size.Y))
		}

		// TODO dont hardcode 4000 as the framebuffer size
		imgui.ImageButtonV(oglToImguiTextureId(w.fb.TextureId()),
			size, //imgui.Vec2{X: wSize.X, Y: wSize.Y},
			imgui.Vec2{X: 0, Y: size.Y / 4000},
			imgui.Vec2{X: size.X / 4000, Y: 0},
			0,
			imgui.Vec4{X: 0, Y: 0, Z: 0, W: 1}, imgui.Vec4{X: 1, Y: 1, Z: 1, W: 1})

	}
	imgui.End()
}
