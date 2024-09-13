package cimgui

// #include <stdlib.h>
// #include <memory.h>
// #include "../cimgui/extra_types.h"
// #include "cimgui_wrapper.h"
// #include "cimgui_typedefs.h"
import "C"

// Export some methods that are necessary for externally packaged backends

type CImTextureID C.ImTextureID

func NewTextureIDFromC(cvalue *CImTextureID) *TextureID {
	return newTextureIDFromC((*C.ImTextureID)(cvalue))
}

func (i Vec4) ToC() C.ImVec4 {
	return i.toC()
}
