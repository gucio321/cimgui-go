package main

import (
	"fmt"
	"github.com/kpango/glg"
	"os"
	"strings"
)

func proceedCallbacks(
	prefix string, callbacks []CIdentifier,
	typedefs *Typedefs, validStructNames []CIdentifier, enums []EnumDef, refTypedefs map[CIdentifier]string,
) (validTypeNames []CIdentifier, err error) {
	validTypeNames = make([]CIdentifier, 0)

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

	callbacksCppSb := &strings.Builder{}
	callbacksCppSb.WriteString(cppFileHeader)
	fmt.Fprintf(callbacksCppSb,
		`
#include "%[1]s_callbacks.h"
#include "cimgui/%[1]s.h"
`, prefix)

	for _, callback := range callbacks {
		if _, ok := refTypedefs[callback]; ok {
			glg.Infof("Skipping %s because it is duplicated.", callback)
			continue
		}

		typedef := typedefs.data[callback]
		glg.Infof("Processing %s (%s)", callback, typedef)
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
