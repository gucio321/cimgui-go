// Code generated by cmd/codegen from https://github.com/AllenDang/cimgui-go.
// DO NOT EDIT.

package ImGuiColorTextEdit

// #include <stdlib.h>
// #include <memory.h>
// #include "../imgui/extra_types.h"
// #include "cimcte_wrapper.h"
// #include "cimcte_typedefs.h"
import "C"
import "github.com/AllenDang/cimgui-go/internal"

type TextEditor struct {
	CData *C.TextEditor
}

// Handle returns C version of TextEditor and its finalizer func.
func (self *TextEditor) Handle() (result *C.TextEditor, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self TextEditor) C() (C.TextEditor, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewTextEditorFromC creates TextEditor from its C pointer.
// SRC ~= *C.TextEditor
func NewTextEditorFromC[SRC any](cvalue SRC) *TextEditor {
	return &TextEditor{CData: internal.ReinterpretCast[*C.TextEditor](cvalue)}
}