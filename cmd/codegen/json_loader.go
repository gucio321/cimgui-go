package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// this stores file bytes of our json files
type jsonData struct {
	structAndEnums,
	typedefs,
	defs,

	refStructAndEnums,
	refTypedefs,

	preset []byte
}

func loadData(f *flags) (*jsonData, error) {
	var err error

	result := &jsonData{}

	result.defs, err = os.ReadFile(f.defJsonPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read definitions json file: %w", err)
	}

	result.typedefs, err = os.ReadFile(f.typedefsJsonpath)
	if err != nil {
		return nil, fmt.Errorf("cannot read typedefs json file: %w", err)
	}

	result.structAndEnums, err = os.ReadFile(f.enumsJsonpath)
	if err != nil {
		log.Panic(err)
	}

	if len(f.refEnumsJsonPath) > 0 {
		result.refStructAndEnums, err = os.ReadFile(f.refEnumsJsonPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read reference struct and enums json file: %w", err)
		}
	}

	if len(f.refTypedefsJsonPath) > 0 {
		result.refTypedefs, err = os.ReadFile(f.refTypedefsJsonPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read reference typedefs json file: %w", err)
		}
	}

	if len(f.presetJsonPath) > 0 {
		result.preset, err = os.ReadFile(f.presetJsonPath)
		if err != nil {
			log.Panic(err)
		}
	}

	return result, nil
}

func (j *jsonData) parseJson() (*Context, error) {
	var err error

	result := &Context{
		arrayIndexGetters: make(map[CIdentifier]CIdentifier),
		preset:            &Preset{},
	}

	if len(j.preset) > 0 {
		json.Unmarshal(j.preset, result.preset)
	}

	// get definitions from json file
	result.funcs, err = getFunDefs(j.defs)
	if err != nil {
		return nil, fmt.Errorf("cannot get function definitions: %w", err)
	}

	result.enums, err = getEnumDefs(j.structAndEnums)
	if err != nil {
		return nil, fmt.Errorf("cannot get enum definitions: %w", err)
	}

	result.typedefs, err = getTypedefs(j.typedefs)
	if err != nil {
		return nil, fmt.Errorf("cannot get typedefs: %w", err)
	}

	result.structs, err = getStructDefs(j.structAndEnums)
	if err != nil {
		return nil, fmt.Errorf("cannot get struct definitions: %w", err)
	}

	result.refTypedefs = make(map[CIdentifier]bool)
	if len(j.refTypedefs) > 0 {
		typedefs, err := getTypedefs(j.refTypedefs)
		if err != nil {
			return nil, fmt.Errorf("cannot get reference typedefs: %w", err)
		}

		result.refTypedefs = RemoveMapValues(typedefs.data)
	}
	_, structs, err := getEnumAndStructNames(j.structAndEnums, result)
	result.structNames = SliceToMap(structs)
	if err != nil {
		return nil, fmt.Errorf("cannot get reference struct and enums names: %w", err)
	}

	result.refEnumNames = make(map[CIdentifier]bool)
	result.refStructNames = make(map[CIdentifier]bool)
	if len(j.refStructAndEnums) > 0 {
		refEnums, refStructs, err := getEnumAndStructNames(j.refStructAndEnums, result)
		if err != nil {
			return nil, fmt.Errorf("cannot get reference struct and enums names: %w", err)
		}

		result.refEnumNames = make(map[CIdentifier]bool)
		for _, refEnum := range refEnums {
			result.refEnumNames[refEnum.renameEnum()] = true
		}

		result.refStructNames = SliceToMap(refStructs)
	}

	return result, nil
}
