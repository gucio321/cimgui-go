// Code generated by cmd/codegen from https://github.com/AllenDang/cimgui-go.
// DO NOT EDIT.

package immarkdown

// #include <stdlib.h>
// #include <memory.h>
// #include "wrapper.h"
// #include "typedefs.h"
// #include "../imgui/extra_types.h"
import "C"
import "github.com/AllenDang/cimgui-go/internal"

type Emphasis struct {
	CData *C.Emphasis
}

// Handle returns C version of Emphasis and its finalizer func.
func (self *Emphasis) Handle() (result *C.Emphasis, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self Emphasis) C() (C.Emphasis, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyEmphasis creates Emphasis with its 0 value.
func NewEmptyEmphasis() *Emphasis {
	return &Emphasis{CData: new(C.Emphasis)}
}

// NewEmphasisFromC creates Emphasis from its C pointer.
// SRC ~= *C.Emphasis
func NewEmphasisFromC[SRC any](cvalue SRC) *Emphasis {
	return &Emphasis{CData: internal.ReinterpretCast[*C.Emphasis](cvalue)}
}

type Line struct {
	CData *C.Line
}

// Handle returns C version of Line and its finalizer func.
func (self *Line) Handle() (result *C.Line, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self Line) C() (C.Line, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyLine creates Line with its 0 value.
func NewEmptyLine() *Line {
	return &Line{CData: new(C.Line)}
}

// NewLineFromC creates Line from its C pointer.
// SRC ~= *C.Line
func NewLineFromC[SRC any](cvalue SRC) *Line {
	return &Line{CData: internal.ReinterpretCast[*C.Line](cvalue)}
}

type Link struct {
	CData *C.Link
}

// Handle returns C version of Link and its finalizer func.
func (self *Link) Handle() (result *C.Link, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self Link) C() (C.Link, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyLink creates Link with its 0 value.
func NewEmptyLink() *Link {
	return &Link{CData: new(C.Link)}
}

// NewLinkFromC creates Link from its C pointer.
// SRC ~= *C.Link
func NewLinkFromC[SRC any](cvalue SRC) *Link {
	return &Link{CData: internal.ReinterpretCast[*C.Link](cvalue)}
}

type MarkdownConfig struct {
	CData *C.MarkdownConfig
}

// Handle returns C version of MarkdownConfig and its finalizer func.
func (self *MarkdownConfig) Handle() (result *C.MarkdownConfig, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self MarkdownConfig) C() (C.MarkdownConfig, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyMarkdownConfig creates MarkdownConfig with its 0 value.
func NewEmptyMarkdownConfig() *MarkdownConfig {
	return &MarkdownConfig{CData: new(C.MarkdownConfig)}
}

// NewMarkdownConfigFromC creates MarkdownConfig from its C pointer.
// SRC ~= *C.MarkdownConfig
func NewMarkdownConfigFromC[SRC any](cvalue SRC) *MarkdownConfig {
	return &MarkdownConfig{CData: internal.ReinterpretCast[*C.MarkdownConfig](cvalue)}
}

type MarkdownFormatInfo struct {
	CData *C.MarkdownFormatInfo
}

// Handle returns C version of MarkdownFormatInfo and its finalizer func.
func (self *MarkdownFormatInfo) Handle() (result *C.MarkdownFormatInfo, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self MarkdownFormatInfo) C() (C.MarkdownFormatInfo, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyMarkdownFormatInfo creates MarkdownFormatInfo with its 0 value.
func NewEmptyMarkdownFormatInfo() *MarkdownFormatInfo {
	return &MarkdownFormatInfo{CData: new(C.MarkdownFormatInfo)}
}

// NewMarkdownFormatInfoFromC creates MarkdownFormatInfo from its C pointer.
// SRC ~= *C.MarkdownFormatInfo
func NewMarkdownFormatInfoFromC[SRC any](cvalue SRC) *MarkdownFormatInfo {
	return &MarkdownFormatInfo{CData: internal.ReinterpretCast[*C.MarkdownFormatInfo](cvalue)}
}

type MarkdownHeadingFormat struct {
	CData *C.MarkdownHeadingFormat
}

// Handle returns C version of MarkdownHeadingFormat and its finalizer func.
func (self *MarkdownHeadingFormat) Handle() (result *C.MarkdownHeadingFormat, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self MarkdownHeadingFormat) C() (C.MarkdownHeadingFormat, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyMarkdownHeadingFormat creates MarkdownHeadingFormat with its 0 value.
func NewEmptyMarkdownHeadingFormat() *MarkdownHeadingFormat {
	return &MarkdownHeadingFormat{CData: new(C.MarkdownHeadingFormat)}
}

// NewMarkdownHeadingFormatFromC creates MarkdownHeadingFormat from its C pointer.
// SRC ~= *C.MarkdownHeadingFormat
func NewMarkdownHeadingFormatFromC[SRC any](cvalue SRC) *MarkdownHeadingFormat {
	return &MarkdownHeadingFormat{CData: internal.ReinterpretCast[*C.MarkdownHeadingFormat](cvalue)}
}

type MarkdownImageData struct {
	CData *C.MarkdownImageData
}

// Handle returns C version of MarkdownImageData and its finalizer func.
func (self *MarkdownImageData) Handle() (result *C.MarkdownImageData, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self MarkdownImageData) C() (C.MarkdownImageData, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyMarkdownImageData creates MarkdownImageData with its 0 value.
func NewEmptyMarkdownImageData() *MarkdownImageData {
	return &MarkdownImageData{CData: new(C.MarkdownImageData)}
}

// NewMarkdownImageDataFromC creates MarkdownImageData from its C pointer.
// SRC ~= *C.MarkdownImageData
func NewMarkdownImageDataFromC[SRC any](cvalue SRC) *MarkdownImageData {
	return &MarkdownImageData{CData: internal.ReinterpretCast[*C.MarkdownImageData](cvalue)}
}

type MarkdownLinkCallbackData struct {
	CData *C.MarkdownLinkCallbackData
}

// Handle returns C version of MarkdownLinkCallbackData and its finalizer func.
func (self *MarkdownLinkCallbackData) Handle() (result *C.MarkdownLinkCallbackData, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self MarkdownLinkCallbackData) C() (C.MarkdownLinkCallbackData, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyMarkdownLinkCallbackData creates MarkdownLinkCallbackData with its 0 value.
func NewEmptyMarkdownLinkCallbackData() *MarkdownLinkCallbackData {
	return &MarkdownLinkCallbackData{CData: new(C.MarkdownLinkCallbackData)}
}

// NewMarkdownLinkCallbackDataFromC creates MarkdownLinkCallbackData from its C pointer.
// SRC ~= *C.MarkdownLinkCallbackData
func NewMarkdownLinkCallbackDataFromC[SRC any](cvalue SRC) *MarkdownLinkCallbackData {
	return &MarkdownLinkCallbackData{CData: internal.ReinterpretCast[*C.MarkdownLinkCallbackData](cvalue)}
}

type MarkdownTooltipCallbackData struct {
	CData *C.MarkdownTooltipCallbackData
}

// Handle returns C version of MarkdownTooltipCallbackData and its finalizer func.
func (self *MarkdownTooltipCallbackData) Handle() (result *C.MarkdownTooltipCallbackData, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self MarkdownTooltipCallbackData) C() (C.MarkdownTooltipCallbackData, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyMarkdownTooltipCallbackData creates MarkdownTooltipCallbackData with its 0 value.
func NewEmptyMarkdownTooltipCallbackData() *MarkdownTooltipCallbackData {
	return &MarkdownTooltipCallbackData{CData: new(C.MarkdownTooltipCallbackData)}
}

// NewMarkdownTooltipCallbackDataFromC creates MarkdownTooltipCallbackData from its C pointer.
// SRC ~= *C.MarkdownTooltipCallbackData
func NewMarkdownTooltipCallbackDataFromC[SRC any](cvalue SRC) *MarkdownTooltipCallbackData {
	return &MarkdownTooltipCallbackData{CData: internal.ReinterpretCast[*C.MarkdownTooltipCallbackData](cvalue)}
}

type TextBlock struct {
	CData *C.TextBlock
}

// Handle returns C version of TextBlock and its finalizer func.
func (self *TextBlock) Handle() (result *C.TextBlock, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self TextBlock) C() (C.TextBlock, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyTextBlock creates TextBlock with its 0 value.
func NewEmptyTextBlock() *TextBlock {
	return &TextBlock{CData: new(C.TextBlock)}
}

// NewTextBlockFromC creates TextBlock from its C pointer.
// SRC ~= *C.TextBlock
func NewTextBlockFromC[SRC any](cvalue SRC) *TextBlock {
	return &TextBlock{CData: internal.ReinterpretCast[*C.TextBlock](cvalue)}
}

type TextRegion struct {
	CData *C.TextRegion
}

// Handle returns C version of TextRegion and its finalizer func.
func (self *TextRegion) Handle() (result *C.TextRegion, fin func()) {
	return self.CData, func() {}
}

// C is like Handle but returns plain type instead of pointer.
func (self TextRegion) C() (C.TextRegion, func()) {
	result, fn := self.Handle()
	return *result, fn
}

// NewEmptyTextRegion creates TextRegion with its 0 value.
func NewEmptyTextRegion() *TextRegion {
	return &TextRegion{CData: new(C.TextRegion)}
}

// NewTextRegionFromC creates TextRegion from its C pointer.
// SRC ~= *C.TextRegion
func NewTextRegionFromC[SRC any](cvalue SRC) *TextRegion {
	return &TextRegion{CData: internal.ReinterpretCast[*C.TextRegion](cvalue)}
}