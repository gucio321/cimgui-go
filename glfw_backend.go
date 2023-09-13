//go:build !exclude_cimgui_glfw
// +build !exclude_cimgui_glfw

package imgui

/*
#define GL_SILENCE_DEPRECATION
#define CIMGUI_USE_GLFW
#define CIMGUI_USE_OPENGL3

#include "cimgui/cimgui.h"
#include "cimgui/cimgui_impl.h"
#include <stdlib.h>
*/
import "C"
import (
	"image"
	"unsafe"

	glfw "github.com/go-gl/glfw/v3.3/glfw"
)

const MAX_EXTRA_FRAME_COUNT = 15
const GLFWTargetFPS = 30

type GLFWWindowFlags int

const (
// TODO
// GLFWWindowFlagsNone         GLFWWindowFlags = GLFWWindowFlags(C.GLFWWindowFlagsNone)
// GLFWWindowFlagsNotResizable GLFWWindowFlags = GLFWWindowFlags(C.GLFWWindowFlagsNotResizable)
// GLFWWindowFlagsMaximized    GLFWWindowFlags = GLFWWindowFlags(C.GLFWWindowFlagsMaximized)
// GLFWWindowFlagsFloating     GLFWWindowFlags = GLFWWindowFlags(C.GLFWWindowFlagsFloating)
// GLFWWindowFlagsFrameless    GLFWWindowFlags = GLFWWindowFlags(C.GLFWWindowFlagsFrameless)
// GLFWWindowFlagsTransparent  GLFWWindowFlags = GLFWWindowFlags(C.GLFWWindowFlagsTransparent)
)

type voidCallbackFunc func()

type GLFWBackend struct {
	afterCreateContext   voidCallbackFunc
	beforeRender         voidCallbackFunc
	afterRender          voidCallbackFunc
	beforeDestoryContext voidCallbackFunc
	window               *glfw.Window

	targetFps       uint
	extraFrameCount int
	bgColor         Vec4
}

func NewGLFWBackend() *GLFWBackend {
	return &GLFWBackend{
		targetFps: GLFWTargetFPS,
		bgColor: Vec4{
			X: 0.45,
			Y: 0.55,
			Z: 0.60,
			W: 1.00,
		},
	}
}

func (b *GLFWBackend) SetAfterCreateContextHook(hook func()) {
	b.afterCreateContext = hook
}

func (b *GLFWBackend) afterCreateHook() func() {
	return b.afterCreateContext
}

func (b *GLFWBackend) SetBeforeDestroyContextHook(hook func()) {
	b.beforeDestoryContext = hook
}

func (b *GLFWBackend) beforeDestroyHook() func() {
	return b.beforeDestoryContext
}

func (b *GLFWBackend) SetBeforeRenderHook(hook func()) {
	b.beforeRender = hook
}

func (b *GLFWBackend) beforeRenderHook() func() {
	return b.beforeRender
}

func (b *GLFWBackend) SetAfterRenderHook(hook func()) {
	b.afterRender = hook
}

func (b *GLFWBackend) afterRenderHook() func() {
	return b.afterRender
}

func (b *GLFWBackend) SetBgColor(color Vec4) {
	C.igSetBgColor(color.toC())
}

func (b *GLFWBackend) Run(loop func()) {
	b.window.MakeContextCurrent()

	// Load Fonts
	// - If no fonts are loaded, dear imgui will use the default font. You can
	// also load multiple fonts and use igPushFont()/PopFont() to select them.
	// - AddFontFromFileTTF() will return the ImFont* so you can store it if you
	// need to select the font among multiple.
	// - If the file cannot be loaded, the function will return NULL. Please
	// handle those errors in your application (e.g. use an assertion, or display
	// an error and quit).
	// - The fonts will be rasterized at a given size (w/ oversampling) and stored
	// into a texture when calling ImFontAtlas::Build()/GetTexDataAsXXXX(), which
	// ImGui_ImplXXXX_NewFrame below will call.
	// - Read 'docs/FONTS.md' for more instructions and details.
	// - Remember that in C/C++ if you want to include a backslash \ in a string
	// literal you need to write a double backslash \\ !
	// io.Fonts->AddFontDefault();
	// io.Fonts->AddFontFromFileTTF("../../misc/fonts/Roboto-Medium.ttf", 16.0f);
	// io.Fonts->AddFontFromFileTTF("../../misc/fonts/Cousine-Regular.ttf", 15.0f);
	// io.Fonts->AddFontFromFileTTF("../../misc/fonts/DroidSans.ttf", 16.0f);
	// io.Fonts->AddFontFromFileTTF("../../misc/fonts/ProggyTiny.ttf", 10.0f);
	// ImFont* font =
	// io.Fonts->AddFontFromFileTTF("c:\\Windows\\Fonts\\ArialUni.ttf", 18.0f,
	// NULL, io.Fonts->GetGlyphRangesJapanese()); IM_ASSERT(font != NULL);

	// Main loop
	lasttime := glfw.GetTime()
	for !b.window.ShouldClose() {
		if b.beforeRender != nil {
			b.beforeRender()
		}

		// render
		// Start the Dear ImGui frame
		C.ImGui_ImplOpenGL3_NewFrame()
		C.ImGui_ImplGlfw_NewFrame()
		NewFrame()

		//b.window.SetUserPointer((unsafe.Pointer)(b.loop))

		// Do ui stuff here
		if loop != nil {
			loop()
		}

		// Rendering
		Render()
		display_w, display_h := b.window.GetFramebufferSize()
		C.glViewport(0, 0, display_w, display_h)
		C.glClearColor(
			b.bgColor.X*b.bgColor.W,
			b.bgColor.Y*b.bgColor.W,
			b.bgColor.Z*b.bgColor.W,
			b.bgColor.W,
		)
		C.glClear(C.GL_COLOR_BUFFER_BIT)
		C.ImGui_ImplOpenGL3_RenderDrawData(CurrentDrawData().handle())

		io := CurrentIO()

		// Update and Render additional Platform Windows
		// (Platform functions may change the current OpenGL context, so we
		// save/restore it to make it easier to paste this code elsewhere.
		//  For this specific demo app we could also call
		//  glfwMakeContextCurrent(window) directly)
		if io.ConfigFlags()&ConfigFlagsViewportsEnable != 0 {
			backup_current_context := glfw.GetCurrentContext()
			UpdatePlatformWindows()
			RenderPlatformWindowsDefault()
			backup_current_context.MakeContextCurrent()
		}

		b.window.SwapBuffers()
		// render end

		for glfw.GetTime() < lasttime+1.0/float64(b.targetFps) {
			// do nothing here
		}
		lasttime += 1.0 / float64(b.targetFps)

		if b.extraFrameCount > 0 {
			b.extraFrameCount--
		} else {
			glfw.WaitEvents()
			b.extraFrameCount = MAX_EXTRA_FRAME_COUNT
		}

		glfw.PollEvents()

		if b.afterRender != nil {
			b.afterRender()
		}
	}

	// Cleanup
	C.ImGui_ImplOpenGL3_Shutdown()
	C.ImGui_ImplGlfw_Shutdown()

	if b.beforeDestroyHook() != nil {
		b.beforeDestroyHook()
	}

	DestroyContext()

	b.window.Destroy()
	glfw.Terminate()
}

func (b *GLFWBackend) SetWindowPos(x, y int) {
	b.window.SetPos(x, y)
}

func (b *GLFWBackend) GetWindowPos() (x, y int32) {
	posX, posY := b.window.GetPos()
	return int32(posX), int32(posY)
}

func (b *GLFWBackend) SetWindowSize(width, height int) {
	b.window.SetSize(width, height)
}

// TODO
func (b GLFWBackend) DisplaySize() (width int32, height int32) {
	widthArg, widthFin := WrapNumberPtr[C.int, int32](&width)
	defer widthFin()

	heightArg, heightFin := WrapNumberPtr[C.int, int32](&height)
	defer heightFin()

	C.igGLFWWindow_GetDisplaySize(b.handle(), widthArg, heightArg)

	return
}

func (b *GLFWBackend) SetWindowTitle(title string) {
	b.window.SetTitle(title)
}

// The minimum and maximum size of the content area of a windowed mode window.
// To specify only a minimum size or only a maximum one, set the other pair to -1
// e.g. SetWindowSizeLimits(640, 480, -1, -1)
func (b *GLFWBackend) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	b.window.SetSizeLimits(minWidth, minHeight, maxWidth, maxHeight)
}

func (b GLFWBackend) SetShouldClose(value bool) {
	b.window.SetShouldClose(value)
}

// TODO: clearify panics
// TODO: fix this gles stuff
func (b *GLFWBackend) CreateWindow(title string, width, height int, flags GLFWWindowFlags) {
	// Setup window
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	// Decide GL+GLSL versions
	//#if defined(IMGUI_IMPL_OPENGL_ES2)
	// GL ES 2.0 + GLSL 100
	//const char *glsl_version = "#version 100";
	//glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 2);
	//glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 0);
	//glfwWindowHint(GLFW_CLIENT_API, GLFW_OPENGL_ES_API);
	//#elif defined(__APPLE__)
	// GL 3.2 + GLSL 150
	//const char *glsl_version = "#version 150";
	//glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
	//glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 2);
	//glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE); // 3.2+ only
	//glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE);           // Required on Mac
	//#else
	// GL 3.0 + GLSL 130
	//const char *glsl_version = "#version 130";
	//glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
	//glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 0);
	//glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);  // 3.2+
	// only glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE); // 3.0+ only
	//#endif

	// Create window with graphics context
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if window == nil || err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1) // Enable vsync

	// Setup Dear ImGui context
	CreateContext()

	if b.afterCreateContext != nil {
		b.afterCreateContext()
	}

	io := CurrentIO()
	igCfgFlags := io.ConfigFlags()
	igCfgFlags |= ConfigFlagsNavEnableKeyboard // Enable Keyboard Controls
	igCfgFlags |= ConfigFlagsDockingEnable     // Enable Docking
	igCfgFlags |= ConfigFlagsViewportsEnable   // Enable Multi-Viewport
	io.SetConfigFlags(igCfgFlags)
	io.SetIniFilename("")

	// Setup Dear ImGui style
	StyleColorsDark()
	// igStyleColorsLight();

	// When viewports are enabled we tweak WindowRounding/WindowBg so platform
	// windows can look identical to regular ones.
	style := CurrentStyle()
	if io.ConfigFlags()&ConfigFlagsViewportsEnable != 0 {
		style.SetWindowRounding(0.0)
		//style->Colors[ImGuiCol_WindowBg].w = 1.0f; // TODO: implement this
	}

	// Setup Platform/Renderer backends
	C.ImGui_ImplGlfw_InitForOpenGL(window, true)
	C.ImGui_ImplOpenGL3_Init(C.glsl_version)

	// Install extra callback
	window.SetRefreshCallback(C.glfw_window_refresh_callback)

	b.window = window
}

func (b *GLFWBackend) SetTargetFPS(fps uint) {
	b.targetFps = fps
}

func (b *GLFWBackend) Refresh() {
	glfw.PostEmptyEvent()
}

// TODO
func (b *GLFWBackend) CreateTexture(pixels unsafe.Pointer, width, height int) TextureID {
	return TextureID(C.igCreateTexture((*C.uchar)(pixels), C.int(width), C.int(height)))
}

// TODO
func (b *GLFWBackend) CreateTextureRgba(img *image.RGBA, width, height int) TextureID {
	return TextureID(C.igCreateTexture((*C.uchar)(&(img.Pix[0])), C.int(width), C.int(height)))
}

// TODO
func (b *GLFWBackend) DeleteTexture(id TextureID) {
	C.igDeleteTexture(C.ImTextureID(id))
}

// SetDropCallback sets the drop callback which is called when an object
// is dropped over the window.
func (b *GLFWBackend) SetDropCallback(cbfun DropCallback) {
	b.window.SetDropCallback(func(w *glfw.Window, names []string) {
		cbfun(names)
	})
}

func (b *GLFWBackend) SetCloseCallback(cbfun WindowCloseCallback) {
	b.window.SetCloseCallback(func(w *glfw.Window) {
		cbfun(b)
	})
}

// SetWindowHint applies to next CreateWindow call
// so use it before CreateWindow call ;-)
func (b *GLFWBackend) SetWindowHint(hint, value int) {
	glfw.WindowHint(glfw.Hint(hint), value) // TODO: abstraction layer - NOT int
}

// SetIcons sets icons for the window.
// THIS CODE COMES FROM https://github.com/go-gl/glfw (BSD-3 clause) - Copyright (c) 2012 The glfw3-go Authors. All rights reserved.
func (b *GLFWBackend) SetIcons(images ...image.Image) {
	b.window.SetIcon(images)
}

func (b *GLFWBackend) SetKeyCallback(cbfun KeyCallback) {
	b.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		cbfun(int(key), scancode, int(action), int(mods)) // TODO: another abstraction layer needed here
	})
}

func (b *GLFWBackend) SetSizeChangeCallback(cbfun SizeChangeCallback) {
	b.window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		cbfun(width, height)
	})
}
