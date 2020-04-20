package fcore

// Platform is the simplified interface to a glfw window
type Platform interface {
	SetCursorEnabled(state bool)
	KeyDown(key rune) bool
}
