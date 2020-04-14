package forgery

import (
	"fmt"

	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/forgery/fcore"
	imguiBackend "github.com/emily33901/forgery/imgui"
	"github.com/emily33901/forgery/loader/keyvalues"
	"github.com/emily33901/forgery/render"
	"github.com/emily33901/forgery/render/adapters"
	"github.com/emily33901/forgery/render/scene"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"github.com/inkyblackness/imgui-go"
)

type Forgery struct {
	core.IDispatcher
	window   window.IWindow
	keyState *window.KeyState
	renderer *renderer.Renderer

	imguiPlatform *imguiBackend.GLFW
	imguiRenderer *imguiBackend.OpenGL3
	context       *imgui.Context

	ShouldQuit bool

	showDemoWindow  bool
	showAboutWindow bool

	Adapter render.Adapter

	fs *filesystem.Filesystem
}

var f *Forgery

// Get gets or creates the forgery singleton
func Get() *Forgery {
	if f != nil {
		return f
	}

	f = &Forgery{
		Adapter:     &adapters.OpenGL{},
		context:     imgui.CreateContext(nil),
		IDispatcher: core.NewDispatcher(),
	}
	f.showDemoWindow = true

	fcore.SetEvents(f.IDispatcher)

	err := window.Init(1280, 720, "Forgery")

	if err != nil {
		panic(err)
	}

	io := imgui.CurrentIO()
	io.SetConfigFlags(imgui.ConfigFlagNavEnableKeyboard)

	f.imguiPlatform = imguiBackend.NewPlatform(io)
	f.imguiRenderer, err = imguiBackend.NewOpenGL3(io)

	if err != nil {
		panic(err)
	}

	f.window = window.Get()
	f.keyState = window.NewKeyState(f.window)
	f.renderer = renderer.NewRenderer(f.window.Gls())
	err = f.renderer.AddDefaultShaders()

	gameinfo, err := keyvalues.FromDisk("E:\\Steam\\steamapps\\common\\Counter-Strike Global Offensive\\csgo\\" + "gameinfo.txt")

	if err != nil {
		panic(err)
	}

	f.fs = filesystem.CreateFromGameDir("E:\\Steam\\steamapps\\common\\Counter-Strike Global Offensive\\csgo\\", gameinfo)

	fmt.Println("Found", len(f.fs.AllPaths()), "Paths")

	if err != nil {
		panic(err)
	}

	return f
}

func (f *Forgery) newSceneWindow() {
	newWindow := scene.NewWindow("", f.Adapter, f.imguiPlatform)

	// Create a blue torus and add it to the scene
	geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)
	newWindow.Scene.Add(mesh)

	// Create and add lights to the scene
	newWindow.Scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	newWindow.Scene.Add(pointLight)

	// Create and add an axis helper to the scene
	newWindow.Scene.Add(helper.NewAxes(0.5))

	newWindow.Camera().SetPosition(0, 0, 3)
}

func (f *Forgery) buildUI() {
	if f.showDemoWindow {
		imgui.ShowDemoWindow(&f.showDemoWindow)
	}

	if f.showAboutWindow {
		f.aboutWindow()
	}

	// Global forgery menu
	if imgui.BeginMainMenuBar() {
		f.menuBar()
		imgui.EndMainMenuBar()
	}

	scene.Iter(func(_ string, v *scene.Window) {
		v.BuildUI(f.imguiPlatform.DeltaTime)
	})
}

func (f *Forgery) aboutWindow() {
	imgui.BeginV("About Forgery", &f.showAboutWindow, imgui.WindowFlagsAlwaysAutoResize|imgui.WindowFlagsNoResize)
	imgui.Text("Forgery - Open source Hammer (r) replacement.")
	imgui.Text("Find the code at https://github.com/emily33901/forgery")
	imgui.Text("Written by Emily Hudson.")
	imgui.Text("Huge thanks to Galaco (https://galaco.me) for the large amounts of code that are his.")
	imgui.End()
}

func (f *Forgery) menuBar() {
	if imgui.MenuItem("New scene") {
		f.newSceneWindow()
	}

	if imgui.BeginMenu("Other") {
		if imgui.MenuItem("About") {
			f.showAboutWindow = true
		}

		if imgui.MenuItemV("Demo window", "", f.showDemoWindow, true) {
			f.showDemoWindow = !f.showDemoWindow
		}

		imgui.EndMenu()
	}
}

func (f *Forgery) render() {
	scene.Iter(func(_ string, v *scene.Window) {
		v.Render(f.renderer)
	})
}

func (f *Forgery) Run() {
	clearColor := [3]float32{0.1, 0.1, 0.1}

	for !f.ShouldQuit && !f.window.(*window.GlfwWindow).ShouldClose() {
		f.imguiPlatform.ProcessEvents()

		f.render()

		f.imguiPlatform.NewFrame()
		imgui.NewFrame()

		f.buildUI()
		imgui.Render()

		f.imguiRenderer.PreRender(clearColor)
		f.imguiRenderer.Render(f.imguiPlatform.DisplaySize(), f.imguiPlatform.FramebufferSize(), imgui.RenderedDrawData())
		f.imguiPlatform.PostRender()
	}

	f.window.Destroy()
}
