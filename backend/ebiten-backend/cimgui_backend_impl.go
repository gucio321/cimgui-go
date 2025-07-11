package ebitenbackend

import (
	"image"
	"image/color"
	"unsafe"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/hajimehoshi/ebiten/v2"
)

func (b *EbitenBackend) SetAfterCreateContextHook(fn func()) {
	b.afterCreateContext = fn
}

func (b *EbitenBackend) SetBeforeDestroyContextHook(fn func()) {
	b.beforeDestroy = fn
}

func (b *EbitenBackend) SetBeforeRenderHook(fn func()) {
	b.beforeRender = fn
}

func (b *EbitenBackend) SetAfterRenderHook(fn func()) {
	b.afterRender = fn
}

func (b *EbitenBackend) SetBgColor(col imgui.Vec4) {
	b.bgColor = col

	fillColor := color.RGBA{
		byte(b.bgColor.X * 255),
		byte(b.bgColor.Y * 255),
		byte(b.bgColor.Z * 255),
		byte(b.bgColor.W * 255),
	}

	toRGBAf32 := func(clr color.Color) (r, g, b, a float32) {
		r16, g16, b16, a16 := clr.RGBA()
		return float32(r16) / 65535.0, float32(g16) / 65535.0, float32(b16) / 65535.0, float32(a16) / 65535.0
	}

	r, g, bC, a := toRGBAf32(fillColor)
	for i := range 4 {
		b.bgColorMagic.pkgFillVertices[i].ColorR = r
		b.bgColorMagic.pkgFillVertices[i].ColorG = g
		b.bgColorMagic.pkgFillVertices[i].ColorB = bC
		b.bgColorMagic.pkgFillVertices[i].ColorA = a
	}
}

func (b *EbitenBackend) Run(loop func()) {
	b.loop = loop
	b.dscale = ebiten.DeviceScaleFactor()

	if b.afterCreateContext != nil {
		b.afterCreateContext()
	}

	if err := ebiten.RunGameWithOptions(b, &ebiten.RunGameOptions{
		ScreenTransparent: b.transparent,
	}); err != nil {
		panic(err)
	}
}

// Because Ebiten refreshes continuously anyway, Refresh has nothing to do here.
func (b *EbitenBackend) Refresh() {
	// noop
}

func (b *EbitenBackend) SetWindowPos(x, y int) {
	ebiten.SetWindowPosition(x, y)
}

func (b *EbitenBackend) GetWindowPos() (x, y int32) {
	xInt, yInt := ebiten.WindowPosition()
	return int32(xInt), int32(yInt)
}

func (b *EbitenBackend) SetWindowSize(width, height int) {
	ebiten.SetWindowSize(width, height)
	/*
		b.width = float32(width)
		b.height = float32(height)
	*/
}

func (b *EbitenBackend) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	ebiten.SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight)
}

func (b *EbitenBackend) SetWindowTitle(title string) {
	ebiten.SetWindowTitle(title)
}

func (b *EbitenBackend) DisplaySize() (width, height int32) {
	w, h := ebiten.Monitor().Size()
	return int32(w), int32(h)
}

// SetShouldClose asks to close the window and stop the game loop. WIll call CloseCallback if set.
func (b *EbitenBackend) SetShouldClose(sc bool) {
	b.shouldClose = sc
}

func (b *EbitenBackend) ContentScale() (xScale, yScale float32) {
	scale := ebiten.DeviceScaleFactor()
	return float32(scale), float32(scale)
}

func (e *EbitenBackend) SetTargetFPS(fps uint) {
	e.fps = fps
}

// TODO: Not implemented
func (b *EbitenBackend) SetDropCallback(backend.DropCallback) {
	panic("SetDropCallback is not implemented for Ebiten backend yet.")
}

func (b *EbitenBackend) SetCloseCallback(cb backend.WindowCloseCallback) {
	b.closeCb = cb
}

func (b *EbitenBackend) SetKeyCallback(backend.KeyCallback) {} // TODO

func (b *EbitenBackend) SetSizeChangeCallback(cb backend.SizeChangeCallback) {
	b.resizeCb = cb
}

func (b *EbitenBackend) SetWindowFlags(flag EbitenBackendFlags, value int) {
	switch flag {
	case EbitenBackendFlagsCursorMode:
		ebiten.SetCursorMode(ebiten.CursorModeType(value))
	case EbitenBackendFlagsCursorShape:
		ebiten.SetCursorShape(ebiten.CursorShapeType(value))
	case EbitenBackendFlagsFPSMode:
		ebiten.SetVsyncEnabled(value == 0)
	case EbitenBackendFlagsResizingMode:
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeType(value))
	case EbitenBackendFlagsDecorated:
		ebiten.SetWindowDecorated(value != 0)
	case EbitenBackendFlagsFloating:
		ebiten.SetWindowFloating(value != 0)
	case EbitenBackendFlagsMaximized:
		if value != 0 {
			ebiten.MaximizeWindow()
		} else {
			ebiten.RestoreWindow()
		}
	case EbitenBackendFlagsMinimized:
		if value != 0 {
			ebiten.MinimizeWindow()
		} else {
			ebiten.RestoreWindow()
		}
	case EbitenBackendFlagsClosingHandled:
		ebiten.SetWindowClosingHandled(value != 0)
	case EbitenBackendFlagsMousePassthrough:
		ebiten.SetWindowMousePassthrough(value != 0)
	case EbitenBackendFlagsDebug:
		b.debug = value != 0
	case EbitenBackendFlagsTransparent:
		b.transparent = value != 0
	default:
		panic("Invalid flag for SetWindowFlags.")
	}
}

func (b *EbitenBackend) SetIcons(icons ...image.Image) {
	ebiten.SetWindowIcon(icons)
}

// TODO: Not implemented
func (b *EbitenBackend) SetSwapInterval(interval EbitenBackendFlags) error {
	panic("SetSwapInterval is not implemented for Ebiten backend yet.")
	return nil
}

// TODO: Not implemented
func (b *EbitenBackend) SetCursorPos(x, y float64) {
	panic("SetCursorPos is not implemented for Ebiten backend yet.")
}

// TODO: Not implemented
func (b *EbitenBackend) SetInputMode(mode, value EbitenBackendFlags) {
	panic("SetInputMode is not implemented for Ebiten backend yet.")
}

func (e *EbitenBackend) CreateWindow(title string, width, height int) {
	// create context
	var imctx *imgui.Context

	switch {
	case e.customCtx != nil:
		imctx = e.customCtx
	case e.customFontAtlas != nil:
		imctx = imgui.CreateContextV(e.customFontAtlas)
	default:
		imctx = imgui.CreateContext()
	}

	e.ctx = imctx

	// indicate that backend (renderer) manages textures.
	// renderer.go will process (load/update/destroy) texture requests.
	io := imgui.CurrentIO()
	io.SetBackendFlags(io.BackendFlags() | imgui.BackendFlagsRendererHasTextures)

	// initialize ebiten stuff
	e.SetWindowTitle(title)
	e.SetWindowSize(width, height)
}

func (e *EbitenBackend) CreateTextureFromGame(game ebiten.Game, width, height int) imgui.TextureRef {
	eimg := ebiten.NewImage(width, height)

	tid := imgui.TextureID(e.cache.NextId())
	e.cache.SetTexture(tid, eimg)
	e.cache.SetGameTexture(tid, game)
	return *imgui.NewTextureRefTextureID(tid)
}

func (e *EbitenBackend) CreateTexture(pixels unsafe.Pointer, width, height int) imgui.TextureRef {
	eimg := ebiten.NewImage(width, height)
	eimg.WritePixels(premultiplyPixels(pixels, width, height))

	tid := imgui.TextureID(e.cache.NextId())
	e.cache.SetTexture(tid, eimg)
	return *imgui.NewTextureRefTextureID(tid)
}

func (b *EbitenBackend) CreateTextureRgba(img *image.RGBA, width, height int) imgui.TextureRef {
	pix := img.Pix
	return b.CreateTexture(unsafe.Pointer(&pix[0]), width, height)
}

func (e *EbitenBackend) DeleteTexture(id imgui.TextureRef) {
	e.cache.RemoveTexture(id.TexID())
}
