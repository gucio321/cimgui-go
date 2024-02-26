package main

import (
	"fmt"
	"github.com/kpango/glg"
	"os"
	"strings"
)

// CallbackArg is a struct when we parse an argument
type CallbackArg struct {
	Type  CIdentifier
	Name  string
	fromC returnWrapper
	toC   ArgumentWrapperData
}

// TODO(gucio321): some wrappers in 2.1 are not needed.
func proceedCallbacks(
	prefix string, callbacks []CIdentifier,
	typedefs *Typedefs, validStructNames []CIdentifier, enums []GoIdentifier, refTypedefs map[CIdentifier]string,
) (validTypeNames map[CIdentifier]bool, err error) {
	validTypeNames = make(map[CIdentifier]bool)

	// 0: Prepare datea
	// 0.1: create writters for files
	// 0.1.1: for GO files
	callbacksGoSb := &strings.Builder{}
	callbacksGoSb.WriteString(goPackageHeader)
	fmt.Fprintf(callbacksGoSb,
		`// #include <stdlib.h>
// #include <memory.h>
// #include "extra_types.h"
// #include "%[1]s_wrapper.h"
// #include "%[1]s_callbacks.h"
import "C"
import "unsafe"
	
var callbackMap = make(map[string]interface{})
`, prefix)

	// 0.1.2: for C Header
	callbacksHeaderSb := &strings.Builder{}
	callbacksHeaderSb.WriteString(cppFileHeader)
	fmt.Fprintf(callbacksHeaderSb,
		`
#pragma once

#include "cimgui/%s.h"

#ifdef __cplusplus
extern "C" {
#endif
`, prefix)

	// 0.1.3: for C file
	callbacksCppSb := &strings.Builder{}
	callbacksCppSb.WriteString(cppFileHeader)
	fmt.Fprintf(callbacksCppSb,
		`
#include "%[1]s_callbacks.h"
#include "cimgui/%[1]s.h"
`, prefix)

	// 0.2: create refeerence data in appropiate format
	validStructNamesMap := make(map[CIdentifier]bool)
	for _, v := range validStructNames {
		validStructNamesMap[v] = true
	}

	enumsMap := make(map[GoIdentifier]bool)
	for _, v := range enums {
		enumsMap[v] = true
	}

callbacksProcess:
	for _, callback := range callbacks {
		if _, ok := refTypedefs[callback]; ok {
			glg.Infof("Skipping %s because it is duplicated.", callback)
			continue
		}

		typedef := typedefs.data[callback]
		glg.Infof("Processing %s (%s)", callback, typedef)

		// 1. Preprocessing: from typedef string extract return type and arguments.
		// 1.1: get return type and arguments list
		// this should be of form something(*)(arg1, arg2)
		// or something name(arg1 name1, arg2 name2) (in immarkdown)
		retArg := Split(typedef, "(*)")
		if len(retArg) != 2 {
			retArg = Split(typedef, fmt.Sprintf(" %s", callback))
			if len(retArg) != 2 {
				glg.Errorf("Callback typedef \"%s\" is of unknown form.", typedef)
				panic("")
			}
		}

		// 1.2: remove redundant brackets and spaces
		ret := retArg[0]
		arg := retArg[1]
		arg = TrimSuffix(arg, ";")
		arg = TrimSuffix(arg, ")")
		arg = TrimPrefix(arg, "(")
		glg.Debugf("-> ret: %s, arg: %s", ret, arg)

		// 1.3: get arguments list
		args := Split(arg, ",")
		argsEx := make([]CallbackArg, len(args))
		for i := range args {
			args[i] = TrimPrefix(args[i], " ")
			args[i] = TrimPrefix(args[i], "const ")
			typeName := Split(args[i], " ")
			ca := CallbackArg{}
			//	Two possibilities:
			//	1. type name
			//	2. type (this also may be "..."
			switch len(typeName) {
			case 1:
				if typeName[0] == "..." {
					continue
				}

				ca.Type = CIdentifier(typeName[0])
				ca.Name = fmt.Sprintf("arg%d", i)
			case 2:
				ca.Type = CIdentifier(typeName[0])
				ca.Name = typeName[1]
			default:
				glg.Errorf("Can't split \"%s\" into type and name part", args[i])
			}

			_, ca.toC, err = getArgWrapper(
				&ArgDef{
					Name: CIdentifier(ca.Name),
					Type: ca.Type,
				},
				false, false,
				validStructNamesMap, map[CIdentifier]bool{}, enumsMap, refTypedefs,
			)

			if err != nil {
				glg.Errorf("Cannot get wrapper for %s: %s", ca.Type, err)
				continue callbacksProcess
			}

			ca.fromC, err = getReturnWrapper(
				ca.Type, validStructNamesMap, enumsMap, refTypedefs,
			)

			if err != nil {
				glg.Errorf("Cannot get wrapper for %s: %s", ca.Type, err)
				continue callbacksProcess
			}

			argsEx[i] = ca
		}

		// 1.4. Find out more about return type
		returnEx := CallbackArg{
			Name: "",
			Type: CIdentifier(ret),
		}

		if returnEx.Type == "void" {
			returnEx.fromC = returnWrapper{
				returnType: "",
				returnStmt: "",
			}
		} else {
			_, returnEx.toC, err = getArgWrapper(
				&ArgDef{
					Name: "",
					Type: returnEx.Type,
				},
				false, false,
				validStructNamesMap, map[CIdentifier]bool{}, enumsMap, refTypedefs,
			)

			if err != nil {
				glg.Errorf("Cannot get wrapper for return type %s: %s", returnEx.Type, err)
				continue callbacksProcess
			}

			returnEx.fromC, err = getReturnWrapper(
				returnEx.Type, validStructNamesMap, enumsMap, refTypedefs,
			)
		}

		// 2. Generate code
		// 2.1. Generate GO code
		goArgs := ""
		for _, a := range argsEx {
			goArgs += fmt.Sprintf("%s %s,", a.Name, a.toC.ArgType)
		}

		goArgs = TrimSuffix(goArgs, ",")

		goCArgs := ""
		for _, a := range argsEx {
			goCArgs += fmt.Sprintf("%s %s,", a.Name, a.toC.CType)
		}

		body := ""

		for _, a := range argsEx {
			body += fmt.Sprintf("%sArg := %s\n", a.Name, fmt.Sprintf(a.fromC.returnStmt, a.Name))
		}

		invocation := ""
		for _, a := range argsEx {
			invocation += fmt.Sprintf("%sArg, ", a.Name)
		}

		returnStmt := fmt.Sprintf("callbackFn(%s)", invocation)
		if returnEx.fromC.returnType == "" {
		}

		fmt.Fprintf(callbacksGoSb,
			`
const mapName_%[1]s = "%[1]s"

type %[1]s func(%[2]s) %[3]s

// export callback%[1]s
func callback%[1]s(%[4]s) %[5]s {
	callbackInterface, ok := callbackMap[mapName_%[1]s]
	if !ok {
		panic("cimgui-go fatal error: Callback %[1]s not found")
	}
	
	callbackFn, ok := callbackInterface.(%[1]s)
	if !ok {
		panic("cimgui-go fatal error: Callback %[1]s is not of proper type")
    }

	%[6]s

	%[7]s
}

func set%[1]sCallback(callback %[1]s) {
	callbackMap[mapName_%[1]s] = callback
}
`, callback.renameGoIdentifier(), goArgs, returnEx.toC.ArgType,
			goCArgs, returnEx.toC.CType,
			body, returnStmt,
		)

		// 2.2. Generate C Header
		fmt.Fprintf(callbacksHeaderSb,
			`
extern %[1]s callback%[2]s(%[3]s);
`, ret, callback.renameGoIdentifier(), Join(args, ", "),
		)

		// 3. Add to valid type names
		validTypeNames[callback] = true
	}

	// 0.3: post processing
	fmt.Fprint(callbacksHeaderSb,
		`
#ifdef __cplusplus
}
#endif`)

	if err := os.WriteFile(fmt.Sprintf("%s_callbacks.go", prefix), []byte(callbacksGoSb.String()), 0644); err != nil {
		return nil, fmt.Errorf("cannot write %s_callbacks.go: %w", prefix, err)
	}

	if err := os.WriteFile(fmt.Sprintf("%s_callbacks.cpp", prefix), []byte(callbacksCppSb.String()), 0644); err != nil {
		return nil, fmt.Errorf("cannot write %s_callbacks.cpp: %w", prefix, err)
	}

	if err := os.WriteFile(fmt.Sprintf("%s_callbacks.h", prefix), []byte(callbacksHeaderSb.String()), 0644); err != nil {
		return nil, fmt.Errorf("cannot write %s_callbacks.h: %w", prefix, err)
	}

	return validTypeNames, nil
}
