package forgery

import (
	"fmt"
	"strings"

	"github.com/emily33901/forgery/core/events"
	"github.com/emily33901/forgery/core/filesystem"
	"github.com/emily33901/forgery/core/materials"
	"github.com/emily33901/forgery/core/textures"
	"github.com/emily33901/forgery/core/vmf"
	"github.com/emily33901/forgery/core/world"
	imguiBackend "github.com/emily33901/forgery/forgery/imgui"
	"github.com/emily33901/forgery/forgery/loader/keyvalues"
	"github.com/emily33901/forgery/forgery/render"
	"github.com/emily33901/forgery/forgery/render/adapters"
	"github.com/emily33901/forgery/forgery/windows"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/renderer"
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

	fs    *filesystem.Filesystem
	world *world.World
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

	events.Set(f.IDispatcher)

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

	for _, f := range f.fs.AllPaths() {
		if strings.HasSuffix(f, ".vtf") && strings.Contains(f, "error") {
			fmt.Println(f)
		}
	}

	if false {
		materials.Load("wow nice meme", f.fs)
		textures.Load("wow nice meme", f.fs)
	}

	vmf, err := vmf.LoadVmf("E:\\emily\\Downloads\\de_60jamey30.vmf")

	if err != nil {
		panic(err)
	}

	f.world = vmf.Worldspawn()

	// textures.Load("wow nice meme", f.fs)

	return f
}

func (f *Forgery) newSceneWindow() {
	newWindow := windows.NewSceneWindow(f.world, "", f.Adapter, f.imguiPlatform)

	// newWindow.Scene.Add(scenes.New())

	// Create and add lights to the scene

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

	windows.Iter(func(_ string, v *windows.SceneWindow) {
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
	windows.Iter(func(_ string, v *windows.SceneWindow) {
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
