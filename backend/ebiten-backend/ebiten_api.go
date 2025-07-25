package ebitenbackend

import (
	"fmt"
	"runtime"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Draw draws the generated imgui frame to the screen.
// This is usually called inside the game's Draw() function.
func (e *EbitenBackend) Draw(screen *ebiten.Image) {
	// add background color
	bounds := screen.Bounds()

	minX, minY := float32(bounds.Min.X), float32(bounds.Min.Y)
	maxX, maxY := float32(bounds.Max.X), float32(bounds.Max.Y)
	e.bgColorMagic.pkgFillVertices[0].DstX = minX
	e.bgColorMagic.pkgFillVertices[0].DstY = minY
	e.bgColorMagic.pkgFillVertices[1].DstX = maxX
	e.bgColorMagic.pkgFillVertices[1].DstY = minY
	e.bgColorMagic.pkgFillVertices[2].DstX = maxX
	e.bgColorMagic.pkgFillVertices[2].DstY = maxY
	e.bgColorMagic.pkgFillVertices[3].DstX = minX
	e.bgColorMagic.pkgFillVertices[3].DstY = maxY
	screen.DrawTriangles(e.bgColorMagic.pkgFillVertices, e.bgColorMagic.pkgFillVertIndices, e.bgColorMagic.pkgMask1x1, &e.bgColorMagic.pkgFillTrianglesOpts)

	if e.debug {
		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("TPS: %.2f\nFPS: %.2f\n[C]lipMask: %t\nSync [I]nput: %t\nSync C[u]rsor: %t\nControl Cursor [S]hape: %t",
				ebiten.ActualTPS(), ebiten.ActualFPS(), e.ClipMask(), e.syncInputs, e.syncCursor, e.controlCursorShape),
			10, 2,
		)
	}

	imgui.Render()
	if e.clipMask {
		if e.lmask == nil {
			w, h := screen.Size()
			e.lmask = ebiten.NewImage(w, h)
		} else {
			w1, h1 := screen.Size()
			w2, h2 := e.lmask.Size()
			if w1 != w2 || h1 != h2 {
				e.lmask.Dispose()
				e.lmask = ebiten.NewImage(w1, h1)
			}
		}
		RenderMasked(screen, e.lmask, imgui.CurrentDrawData(), e.cache, e.filter)
	} else {
		Render(screen, imgui.CurrentDrawData(), e.cache, e.filter)
	}

	// render texture cache game textures
	e.cache.ForEachGame(func(is imgui.TextureID, game ebiten.Game, target *ebiten.Image) {
		target.Clear()
		game.Draw(target)
	})
}

// Update needs to be called on every frame, before cimgui-go calls.
// This is usually called inside the game's Update() function.
// delta is the time in seconds since the last frame.
func (e *EbitenBackend) Update() error {
	e.BeginFrame()
	e.loop()
	e.EndFrame()

	// render texture cache game textures
	e.cache.ForEachGame(func(_ imgui.TextureID, game ebiten.Game, _ *ebiten.Image) {
		game.Update()
	})

	if e.shouldClose {
		if e.closeCb != nil {
			e.closeCb()
		}

		return ebiten.Termination
	}

	return nil
}

func (e *EbitenBackend) Layout(outsideWidth, outsideHeight int) (int, int) {
	var newW, newH int
	if e.retina {
		m := ebiten.DeviceScaleFactor()
		newW = int(float64(outsideWidth) * m)
		newH = int(float64(outsideHeight) * m)
	} else {
		newW = outsideWidth
		newH = outsideHeight
	}

	if e.currentWidth != newW || e.currentHeight != newH {
		e.currentWidth = newW
		e.currentHeight = newH

		if e.resizeCb != nil {
			e.resizeCb(e.currentWidth, e.currentHeight)
		}
	}

	return e.currentWidth, e.currentHeight
}

func (e *EbitenBackend) onfinalize() {
	if e.beforeDestroy != nil {
		e.beforeDestroy()
	}

	runtime.SetFinalizer(e, nil)
	imgui.DestroyContext()
}

func (e *EbitenBackend) controlCursorShapeFn() {
	if !e.controlCursorShape {
		return
	}

	switch imgui.CurrentMouseCursor() {
	case imgui.MouseCursorNone:
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	case imgui.MouseCursorArrow:
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	case imgui.MouseCursorTextInput:
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	case imgui.MouseCursorResizeAll:
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
	case imgui.MouseCursorResizeEW:
		ebiten.SetCursorShape(ebiten.CursorShapeEWResize)
	case imgui.MouseCursorResizeNS:
		ebiten.SetCursorShape(ebiten.CursorShapeNSResize)
	case imgui.MouseCursorHand:
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	default:
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
}
