package imgui

import "C"
import (
	"image"
	"unsafe"
)

var currentBackend Backend

type DropCallback func([]string)
type KeyCallback func(key, scanCode, action, mods int)
type SizeChangeCallback func(w, h int)

type WindowCloseCallback func(b Backend)

// Backend is a special interface that implements all methods required
// to render imgui application.
type Backend interface {
	SetAfterCreateContextHook(func())
	SetBeforeDestroyContextHook(func())
	SetBeforeRenderHook(func())
	SetAfterRenderHook(func())

	SetBgColor(color Vec4)
	Run(func())
	Refresh()

	SetWindowPos(x, y int)
	GetWindowPos() (x, y int32)
	SetWindowSize(width, height int)
	SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int)
	SetWindowTitle(title string)
	DisplaySize() (width, height int32)
	SetShouldClose(bool)

	SetTargetFPS(fps uint)

	CreateTexture(pixels unsafe.Pointer, width, Height int) TextureID
	CreateTextureRgba(img *image.RGBA, width, height int) TextureID
	DeleteTexture(id TextureID)
	SetDropCallback(DropCallback)
	SetCloseCallback(WindowCloseCallback)
	SetKeyCallback(KeyCallback)
	SetSizeChangeCallback(SizeChangeCallback)
	// SetWindowHint selected hint to specified value.
	// For list of hints check GLFW source code.
	// TODO: this needs generic layer
	SetWindowHint(hint, value int)
	SetIcons(icons ...image.Image)

	// TODO: flags needs generic layer
	CreateWindow(title string, width, height int, flags GLFWWindowFlags)
}

func CreateBackend(backend Backend) Backend {
	currentBackend = backend
	return currentBackend
}

func GetBackend() Backend {
	return currentBackend
}
