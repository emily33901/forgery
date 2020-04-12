package imgui

import (
	"math"

	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/inkyblackness/imgui-go"
)

// Implements platform with g3n as a backend
type GLFW struct {
	imguiIO imgui.IO

	window window.IWindow

	time             float64
	mouseJustPressed [3]bool
}

func NewPlatform(io imgui.IO) *GLFW {
	platform := &GLFW{
		imguiIO: io,
		window:  window.Get(),
	}

	platform.setKeyMapping()
	platform.subscribeToEvents()

	return platform
}

func (platform *GLFW) Dispose() {
	platform.window.Destroy()
}

// DisplaySize returns the dimension of the display.
func (platform *GLFW) DisplaySize() [2]float32 {
	w, h := platform.window.GetSize()
	return [2]float32{float32(w), float32(h)}
}

// FramebufferSize returns the dimension of the framebuffer.
func (platform *GLFW) FramebufferSize() [2]float32 {
	w, h := platform.window.GetFramebufferSize()
	return [2]float32{float32(w), float32(h)}
}

func (platform *GLFW) PostRender() {
	platform.window.(*window.GlfwWindow).SwapBuffers()
}

func (platform *GLFW) NewFrame() {
	// Setup display size (every frame to accommodate for window resizing)
	displaySize := platform.DisplaySize()
	platform.imguiIO.SetDisplaySize(imgui.Vec2{X: displaySize[0], Y: displaySize[1]})

	// Setup time step
	currentTime := glfw.GetTime()
	if platform.time > 0 {
		platform.imguiIO.SetDeltaTime(float32(currentTime - platform.time))
	}
	platform.time = currentTime

	// Setup inputs
	if platform.window.(*window.GlfwWindow).GetAttrib(glfw.Focused) != 0 {
		x, y := platform.window.(*window.GlfwWindow).GetCursorPos()
		platform.imguiIO.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})
	} else {
		platform.imguiIO.SetMousePosition(imgui.Vec2{X: -math.MaxFloat32, Y: -math.MaxFloat32})
	}

	for i := 0; i < len(platform.mouseJustPressed); i++ {
		down := platform.mouseJustPressed[i] || (platform.window.(*window.GlfwWindow).GetMouseButton(glfwButtonIDByIndex[i]) == glfw.Press)
		platform.imguiIO.SetMouseButtonDown(i, down)
		platform.mouseJustPressed[i] = false
	}

}

func (platform *GLFW) ProcessEvents() {
	platform.window.(*window.GlfwWindow).PollEvents()
}

func (platform *GLFW) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	platform.imguiIO.KeyMap(imgui.KeyTab, int(glfw.KeyTab))
	platform.imguiIO.KeyMap(imgui.KeyLeftArrow, int(glfw.KeyLeft))
	platform.imguiIO.KeyMap(imgui.KeyRightArrow, int(glfw.KeyRight))
	platform.imguiIO.KeyMap(imgui.KeyUpArrow, int(glfw.KeyUp))
	platform.imguiIO.KeyMap(imgui.KeyDownArrow, int(glfw.KeyDown))
	platform.imguiIO.KeyMap(imgui.KeyPageUp, int(glfw.KeyPageUp))
	platform.imguiIO.KeyMap(imgui.KeyPageDown, int(glfw.KeyPageDown))
	platform.imguiIO.KeyMap(imgui.KeyHome, int(glfw.KeyHome))
	platform.imguiIO.KeyMap(imgui.KeyEnd, int(glfw.KeyEnd))
	platform.imguiIO.KeyMap(imgui.KeyInsert, int(glfw.KeyInsert))
	platform.imguiIO.KeyMap(imgui.KeyDelete, int(glfw.KeyDelete))
	platform.imguiIO.KeyMap(imgui.KeyBackspace, int(glfw.KeyBackspace))
	platform.imguiIO.KeyMap(imgui.KeySpace, int(glfw.KeySpace))
	platform.imguiIO.KeyMap(imgui.KeyEnter, int(glfw.KeyEnter))
	platform.imguiIO.KeyMap(imgui.KeyEscape, int(glfw.KeyEscape))
	platform.imguiIO.KeyMap(imgui.KeyA, int(glfw.KeyA))
	platform.imguiIO.KeyMap(imgui.KeyC, int(glfw.KeyC))
	platform.imguiIO.KeyMap(imgui.KeyV, int(glfw.KeyV))
	platform.imguiIO.KeyMap(imgui.KeyX, int(glfw.KeyX))
	platform.imguiIO.KeyMap(imgui.KeyY, int(glfw.KeyY))
	platform.imguiIO.KeyMap(imgui.KeyZ, int(glfw.KeyZ))
}

func (platform *GLFW) subscribeToEvents() {
	onKeyPress := func(down bool, evname string, ev interface{}) {
		e := ev.(*window.KeyEvent)

		if down {
			platform.imguiIO.KeyPress(int(e.Key))
		} else {
			platform.imguiIO.KeyRelease(int(e.Key))
		}

		platform.imguiIO.KeyCtrl(int(glfw.KeyLeftControl), int(glfw.KeyRightControl))
		platform.imguiIO.KeyShift(int(glfw.KeyLeftShift), int(glfw.KeyRightShift))
		platform.imguiIO.KeyAlt(int(glfw.KeyLeftAlt), int(glfw.KeyRightAlt))
		platform.imguiIO.KeySuper(int(glfw.KeyLeftSuper), int(glfw.KeyRightSuper))
	}

	platform.window.Subscribe(window.OnKeyDown, func(name string, ev interface{}) {
		onKeyPress(true, name, ev)
	})
	platform.window.Subscribe(window.OnKeyUp, func(name string, ev interface{}) {
		onKeyPress(false, name, ev)
	})

	platform.window.Subscribe(window.OnChar, func(evname string, ev interface{}) {
		imgui.CurrentIO().AddInputCharacters(string(ev.(*window.CharEvent).Char))
	})

	platform.window.Subscribe(window.OnScroll, func(evname string, ev interface{}) {
		e := ev.(*window.ScrollEvent)
		platform.imguiIO.AddMouseWheelDelta(float32(e.Xoffset), float32(e.Yoffset))
	})

	platform.window.Subscribe(window.OnMouseDown, func(evname string, ev interface{}) {
		e := ev.(*window.MouseEvent)

		buttonIndex, known := glfwButtonIndexByID[e.Button]

		if known {
			platform.mouseJustPressed[buttonIndex] = true
		}
	})
}

var glfwButtonIndexByID = map[window.MouseButton]int{
	window.MouseButton1: 0,
	window.MouseButton2: 1,
	window.MouseButton3: 2,
}

var glfwButtonIDByIndex = map[int]glfw.MouseButton{
	0: glfw.MouseButton1,
	1: glfw.MouseButton2,
	2: glfw.MouseButton3,
}

// ClipboardText returns the current clipboard text, if available.
func (platform *GLFW) ClipboardText() (string, error) {
	return platform.window.(*window.GlfwWindow).GetClipboardString()
}

// SetClipboardText sets the text as the current clipboard text.
func (platform *GLFW) SetClipboardText(text string) {
	platform.window.(*window.GlfwWindow).SetClipboardString(text)
}
