// Package goglfwbackend aims to provide a GLFW backend for cimgui-go using go-gl/glfw.
package goglfwbackend

import (
	"image"
	"log"
	"unsafe"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type GLFWWindowFlags int

// GLFWBackend implements glfw backend using go-gl/glfw
type GLFWBackend struct {
	afterCreateContext,
	beforeRender,
	loop,
	afterRender,
	beforeDestroyContext func()

	bgColor imgui.Vec4

	window *glfw.Window
}

// NewGLFWBackend creates a new GLFW backend instance
// NOTE: this also initializes go-gl/glfw so make sure no conflict (e.g. with glfwbackend) happens.
func NewGLFWBackend() *GLFWBackend {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	return &GLFWBackend{
		afterCreateContext:   func() {},
		beforeRender:         func() {},
		loop:                 func() {},
		afterRender:          func() {},
		beforeDestroyContext: func() {},
	}
}

func (G GLFWBackend) SetAfterCreateContextHook(f func()) {
	G.afterCreateContext = f
}

func (G GLFWBackend) SetBeforeDestroyContextHook(f func()) {
	G.beforeDestroyContext = f
}

func (G GLFWBackend) SetBeforeRenderHook(f func()) {
	G.beforeRender = f
}

func (G GLFWBackend) SetAfterRenderHook(f func()) {
	G.afterRender = f
}

func (G GLFWBackend) SetBgColor(color imgui.Vec4) {
	G.bgColor = color
}

func (G GLFWBackend) Run(f func()) {
	G.window.MakeContextCurrent()
}

func (G GLFWBackend) Refresh() {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetWindowPos(x, y int) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) GetWindowPos() (x, y int32) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetWindowSize(width, height int) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetWindowTitle(title string) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) DisplaySize() (width, height int32) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetShouldClose(b bool) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) ContentScale() (xScale, yScale float32) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetTargetFPS(fps uint) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetDropCallback(callback backend.DropCallback) {
	//TODO implement me
	log.Print("implement me")
}

func (G GLFWBackend) SetCloseCallback(callback backend.WindowCloseCallback) {
	//TODO implement me
	log.Print("implement me")
}

func (G GLFWBackend) SetKeyCallback(callback backend.KeyCallback) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetSizeChangeCallback(callback backend.SizeChangeCallback) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetWindowFlags(flag GLFWWindowFlags, value int) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetIcons(icons ...image.Image) {
	G.window.SetIcon(icons)
}

func (G GLFWBackend) SetSwapInterval(interval GLFWWindowFlags) error {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetCursorPos(x, y float64) {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) SetInputMode(mode GLFWWindowFlags, value GLFWWindowFlags) {
	//TODO implement me
	panic("implement me")
}

func (G *GLFWBackend) CreateWindow(title string, width, height int) {
	var err error
	G.window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)

	}

	if 
}

func (G GLFWBackend) CreateTexture(pixels unsafe.Pointer, width, Height int) imgui.TextureID {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) CreateTextureRgba(img *image.RGBA, width, height int) imgui.TextureID {
	//TODO implement me
	panic("implement me")
}

func (G GLFWBackend) DeleteTexture(id imgui.TextureID) {
	//TODO implement me
	panic("implement me")
}
