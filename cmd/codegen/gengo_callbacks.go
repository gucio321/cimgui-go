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

func proceedCallbacks(
	prefix string, callbacks []CIdentifier,
	typedefs *Typedefs, validStructNames []CIdentifier, enums []GoIdentifier, refTypedefs map[CIdentifier]string,
) (validTypeNames []CIdentifier, err error) {
	validTypeNames = make([]CIdentifier, 0)

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
//import "C"
//import "unsafe"
	
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
				validStructNamesMap, enumsMap, refTypedefs,
			)

			if err != nil {
				glg.Errorf("Cannot get wrapper for %s: %s", ca.Type, err)
			}

			ca.fromC, err = getReturnWrapper(
				ca.Type, validStructNamesMap, enumsMap, refTypedefs,
			)

			fmt.Println(ca)
			argsEx[i] = ca
		}
	}

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