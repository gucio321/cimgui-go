package main

import "C"
import (
	"fmt"
)

type ArgumentWrapperData struct {
	// go-valid argument type (e.g. string, ImVec2, etc.)
	ArgType string
	// argument deffinition (e.g. arg1, arg1Fin := ...\ndefer arg1Fin())
	ArgDef string
	// one-line go statement that should be called after calling C function
	// in order to update go value.
	// It is intended to be run in defer statement (and it will be in most cases)
	// so it should be one-line function call OR a call of custom function
	Finalizer string

	// a name of variable of wrapped C type
	VarName string
}

type argumentWrapper func(argName string) ArgumentWrapperData

func argWrapper(argType string) (wrapper argumentWrapper, err error) {
	argWrapperMap := map[string]argumentWrapper{
		"char":                     simpleW("rune", "C.char"),
		"char[5]":                  simplePtrArrayW(5, "C.char", "rune"),
		"char*":                    constCharW,
		"const char*":              constCharW,
		"const char**":             charPtrPtrW,
		"const char* const[]":      charPtrPtrW,
		"unsigned char":            simpleW("uint", "C.uchar"),
		"unsigned char*":           simplePtrW("uint", "C.uchar"),
		"unsigned char**":          uCharPtrW,
		"size_t":                   simpleW("uint64", "C.xulong"),
		"size_t*":                  sizeTPtrW,
		"float":                    simpleW("float32", "C.float"),
		"const float":              simpleW("float32", "C.float"),
		"float*":                   floatPtrW,
		"const float*":             floatArrayW,
		"short":                    simpleW("int", "C.short"),
		"unsigned short":           simpleW("uint", "C.ushort"),
		"ImU8":                     simpleW("uint", "C.ImU8"),
		"const ImU8*":              simplePtrSliceW("C.ImU8", "byte"),
		"ImU16":                    simpleW("uint", "C.ImU16"),
		"const ImU16*":             simplePtrSliceW("C.ImU16", "uint16"),
		"ImU32":                    simpleW("uint32", "C.ImU32"),
		"ImU32*":                   simplePtrSliceW("C.ImU32", "uint32"),
		"const ImU32*":             simplePtrSliceW("C.ImU32", "uint32"),
		"ImU64":                    simpleW("uint64", "C.ImU64"),
		"ImU64*":                   simplePtrSliceW("C.ImU64", "uint64"),
		"const ImU64*":             uint64ArrayW,
		"ImS8":                     simpleW("int", "C.ImS8"),
		"const ImS8*":              simplePtrSliceW("C.ImS8", "int8"),
		"ImS16":                    simpleW("int", "C.ImS16"),
		"const ImS16*":             simplePtrSliceW("C.ImS16", "int"),
		"ImS32":                    simpleW("int", "C.ImS32"),
		"const ImS32*":             simplePtrSliceW("C.ImS32", "int32"),
		"ImS64":                    simpleW("int64", "C.ImS64"),
		"ImS64*":                   simplePtrW("int64", "C.ImS64"),
		"const ImS64*":             int64ArrayW,
		"int":                      simpleW("int32", "C.int"),
		"int*":                     simplePtrW("int32", "C.int"),
		"unsigned int":             simpleW("uint32", "C.uint"),
		"unsigned int*":            simplePtrW("uint32", "C.uint"),
		"double":                   simpleW("float64", "C.double"),
		"double*":                  simplePtrW("float64", "C.double"),
		"const double*":            simplePtrSliceW("C.double", "float64"),
		"bool":                     simpleW("bool", "C.bool"),
		"bool*":                    boolPtrW,
		"int[2]":                   simplePtrArrayW(2, "C.int", "int32"),
		"int[3]":                   simplePtrArrayW(3, "C.int", "int32"),
		"int[4]":                   simplePtrArrayW(4, "C.int", "int32"),
		"float[2]":                 simplePtrArrayW(2, "C.float", "float32"),
		"float[3]":                 simplePtrArrayW(3, "C.float", "float32"),
		"float[4]":                 simplePtrArrayW(4, "C.float", "float32"),
		"ImWchar":                  simpleW("Wchar", "C.ImWchar"),
		"ImWchar*":                 simpleW("*Wchar", "(*C.ImWchar)"),
		"const ImWchar*":           simpleW("*Wchar", "(*C.ImWchar)"),
		"ImWchar16":                simpleW("uint16", "C.ImWchar16"),
		"ImGuiID":                  simpleW("ID", "C.ImGuiID"),
		"ImGuiID*":                 simplePtrW("ID", "C.ImGuiID"),
		"ImTextureID":              simpleW("TextureID", "C.ImTextureID"),
		"ImDrawIdx":                simpleW("DrawIdx", "C.ImDrawIdx"),
		"ImGuiTableColumnIdx":      simpleW("TableColumnIdx", "C.ImGuiTableColumnIdx"),
		"ImGuiTableDrawChannelIdx": simpleW("TableDrawChannelIdx", "C.ImGuiTableDrawChannelIdx"),
		"ImGuiKeyChord":            simpleW("KeyChord", "C.ImGuiKeyChord"),
		"void*":                    simpleW("unsafe.Pointer", ""),
		"const void*":              simpleW("unsafe.Pointer", ""),
		"const ImVec2":             wrappableW("Vec2"),
		"const ImVec2*":            wrappablePtrW("*Vec2", "C.ImVec2"),
		"ImVec2":                   wrappableW("Vec2"),
		"ImVec2*":                  wrappablePtrW("*Vec2", "C.ImVec2"),
		"ImVec2[2]":                wrappablePtrArrayW(2, "C.ImVec2", "Vec2"),
		"const ImVec4":             wrappableW("Vec4"),
		"const ImVec4*":            wrappablePtrW("*Vec4", "C.ImVec4"),
		"ImVec4":                   wrappableW("Vec4"),
		"ImVec4*":                  wrappablePtrW("*Vec4", "C.ImVec4"),
		"ImColor*":                 wrappablePtrW("*Color", "C.ImColor"),
		"ImRect":                   wrappableW("Rect"),
		"const ImRect":             wrappableW("Rect"),
		"ImRect*":                  wrappablePtrW("*Rect", "C.ImRect"),
		"const ImRect*":            wrappablePtrW("*Rect", "C.ImRect"),
		"ImPlotPoint":              wrappableW("PlotPoint"),
		"const ImPlotPoint":        wrappableW("PlotPoint"),
		"ImPlotPoint*":             wrappablePtrW("*PlotPoint", "C.ImPlotPoint"),
		"ImPlotTime":               wrappableW("PlotTime"),
		"const ImPlotTime":         wrappableW("PlotTime"),
		"ImPlotTime*":              wrappablePtrW("*PlotTime", "C.ImPlotTime"),
	}

	if wrapper, ok := argWrapperMap[argType]; ok {
		return wrapper, nil
	}

	return nil, fmt.Errorf("no wrapper for type %s", argType)
}

func constCharW(arg string) ArgumentWrapperData {
	return ArgumentWrapperData{
		ArgType:   "string",
		VarName:   fmt.Sprintf("%sArg", arg),
		ArgDef:    fmt.Sprintf("%[1]sArg, %[1]sFin := WrapString(%[1]s)", arg),
		Finalizer: fmt.Sprintf("%sFin()", arg),
	}
}

func charPtrPtrW(arg string) ArgumentWrapperData {
	return ArgumentWrapperData{
		ArgType:   "[]string",
		VarName:   fmt.Sprintf("%sArg", arg),
		ArgDef:    fmt.Sprintf("%[1]sArg, %[1]sFin := WrapStringList(%[1]s)", arg),
		Finalizer: fmt.Sprintf("%sFin()", arg),
	}
}

func uCharPtrW(arg string) ArgumentWrapperData {
	return ArgumentWrapperData{
		ArgType: "*C.uchar",
		VarName: fmt.Sprintf("&%s", arg),
	}
}

func sizeTPtrW(arg string) ArgumentWrapperData {
	return ArgumentWrapperData{
		ArgType: "*uint64",
		VarName: fmt.Sprintf("(*C.xulong)(%s)", arg),
	}
}

// leave this for now because of https://github.com/AllenDang/cimgui-go/issues/31
func floatPtrW(arg string) ArgumentWrapperData {
	return simplePtrW("float32", "C.float")(arg)
}

func floatArrayW(arg string) ArgumentWrapperData {
	return ArgumentWrapperData{
		ArgType: "[]float32",
		VarName: fmt.Sprintf("(*C.float)(&(%s[0]))", arg),
	}
}

func boolPtrW(arg string) ArgumentWrapperData {
	return ArgumentWrapperData{
		ArgType:   "*bool",
		ArgDef:    fmt.Sprintf("%[1]sArg, %[1]sFin := WrapBool(%[1]s)", arg),
		Finalizer: fmt.Sprintf("%[1]sFin()", arg),
		VarName:   fmt.Sprintf("%sArg", arg),
	}
}

func int64ArrayW(arg string) ArgumentWrapperData {
	return ArgumentWrapperData{
		ArgType: "[]int64",
		VarName: fmt.Sprintf("(*C.longlong)(&(%s[0]))", arg),
	}
}

func uint64ArrayW(arg string) ArgumentWrapperData {
	return ArgumentWrapperData{
		ArgType: "[]uint64",
		VarName: fmt.Sprintf("(*C.ulonglong)(&(%s[0]))", arg),
	}
}

// generic wrappers:

// C.int -> int32
func simpleW(goType, cType string) argumentWrapper {
	return func(arg string) ArgumentWrapperData {
		return ArgumentWrapperData{
			ArgType: goType,
			VarName: fmt.Sprintf("%s(%s)", cType, arg),
		}
	}
}

// C.int* -> *int32
//
//	return simplePtrW(arg.Name, "int16", "C.int")
func simplePtrW(goType, cType string) argumentWrapper {
	return func(arg string) ArgumentWrapperData {
		return ArgumentWrapperData{
			ArgType:   fmt.Sprintf("*%s", goType),
			ArgDef:    fmt.Sprintf("%[1]sArg, %[1]sFin := WrapNumberPtr[%[2]s, %[3]s](%[1]s)", arg, cType, goType),
			Finalizer: fmt.Sprintf("%[1]sFin()", arg, cType, goType),
			VarName:   fmt.Sprintf("%sArg", arg),
		}
	}
}

func vectorW(baseGoType string) argumentWrapper {
	return func(arg ArgDef) ArgumentWrapperData {
		return ArgumentWrapperData{
			ArgType: fmt.Sprintf("Vector[%s]", baseGoType),
			ArgDef:  fmt.Sprintf(`%[1]sArg := newVectorFromC(%[1]s.Size, %[1]s.Capacity, new%[2]sFromC(%[1]s.Data))`, arg.Name, baseGoType),
			VarName: fmt.Sprintf("%sArg", arg.Name),
		}
	}
}

// C.int*, C.int[] as well as C.int[2] -> [2]*int32
func simplePtrArrayW(size int, cArrayType, goArrayType string) argumentWrapper {
	return func(arg string) ArgumentWrapperData {
		return ArgumentWrapperData{
			ArgType: fmt.Sprintf("*[%d]%s", size, goArrayType),
			ArgDef: fmt.Sprintf(`
%[1]sArg := make([]%[2]s, len(%[1]s))
for i, %[1]sV := range %[1]s {
  %[1]sArg[i] = %[2]s(%[1]sV)
}`, arg, cArrayType),
			VarName: fmt.Sprintf("(*%s)(&%sArg[0])", cArrayType, arg),
			Finalizer: fmt.Sprintf(`
for i, %[1]sV := range %[1]sArg {
	(*%[1]s)[i] = %[3]s(%[1]sV)
}

`, arg, cArrayType, goArrayType),
		}
	}
}

// C.int*, C.int[] -> *[]int32
func simplePtrSliceW(cArrayType, goArrayType string) argumentWrapper {
	return func(arg string) ArgumentWrapperData {
		return ArgumentWrapperData{
			ArgType: fmt.Sprintf("*[]%s", goArrayType),
			ArgDef: fmt.Sprintf(`%[1]sArg := make([]%[2]s, len(*%[1]s))
for i, %[1]sV := range *%[1]s {
  %[1]sArg[i] = %[2]s(%[1]sV)
}
`, arg, cArrayType, goArrayType),
			Finalizer: fmt.Sprintf(`
  for i, %[1]sV := range %[1]sArg {
    (*%[1]s)[i] = %[3]s(%[1]sV)
  }
`, arg, cArrayType, goArrayType),
			VarName: fmt.Sprintf("(*%s)(&%sArg[0])", cArrayType, arg),
		}
	}
}

// C.ImVec2 -> ImVec2
func wrappableW(sType string) argumentWrapper {
	return func(arg string) ArgumentWrapperData {
		return ArgumentWrapperData{
			ArgType: sType,
			VarName: fmt.Sprintf("%s.toC()", arg),
		}
	}
}

// C.ImVec2* -> *ImVec2
func wrappablePtrW(goType, cType string) argumentWrapper {
	return func(arg string) ArgumentWrapperData {
		return ArgumentWrapperData{
			ArgType:   goType,
			ArgDef:    fmt.Sprintf("%[1]sArg, %[1]sFin := wrap[%[3]s, %[2]s](%[1]s)", arg, goType, cType),
			Finalizer: fmt.Sprintf("%[1]sFin()", arg, goType, cType),
			VarName:   fmt.Sprintf("%sArg", arg),
		}
	}
}

func wrappablePtrArrayW(size int, cArrayType, goArrayType string) argumentWrapper {
	return func(arg string) ArgumentWrapperData {
		return ArgumentWrapperData{
			ArgType: fmt.Sprintf("[%d]*%s", size, goArrayType),
			ArgDef: fmt.Sprintf(`%[1]sArg := make([]%[2]s, len(%[1]s))
%[1]sFin := make([]func(), len(%[1]s))
for i, %[1]sV := range %[1]s {
	var tmp *%[2]s
  	tmp, %[1]sFin[i] = wrap[%[2]s, *%[3]s](%[1]sV)
  	%[1]sArg[i] = *tmp
}
`, arg, cArrayType, goArrayType),
			Finalizer: fmt.Sprintf(`
  for _, %[1]sV := range %[1]sFin {
    %[1]sV()
  }
`, arg, cArrayType, goArrayType),
			VarName: fmt.Sprintf("(*%s)(&%sArg[0])", cArrayType, arg),
		}
	}
}
