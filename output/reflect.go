package output

import (
	"reflect"
	"strings"
)

func ReflectStruct(aStruct any) []Meta {
	var results []Meta
	sc := reflect.TypeOf(aStruct)
	for i := 0; i < sc.NumField(); i++ {
		var meta Meta
		meta.Index = i
		field := sc.Field(i)
		meta.Name = field.Name
		parts := strings.Split(field.Tag.Get("name"), ",")
		if len(parts) > 1 {
			meta.Tag = parts[0]
			meta.Dtype = parts[1]
			meta.Cols = 1
			results = append(results, meta)
		}
	}
	return results
}

func FindActiveScriptCols(structs []Script, meta []Meta) []Meta {
	var results []Meta
	for _, mt := range meta {
		for _, scr := range structs {
			data := reflect.ValueOf(scr).Field(mt.Index)
			if !data.IsZero() {
				results = append(results, mt)
				break
			}
		}
	}
	return results
}

func FindActiveWordCols(structs []Word, meta []Meta) []Meta {
	var results []Meta
	for _, mt := range meta {
		for _, scr := range structs {
			data := reflect.ValueOf(scr).Field(mt.Index)
			if !data.IsZero() {
				results = append(results, mt)
				break
			}
		}
	}
	return results
}

func FindNumScriptMFCC(scripts []Script) int {
	var result int
	for _, scr := range scripts {
		cols := scr.MFCCCols
		if cols > 0 {
			result = cols
			break
		}
	}
	return result
}

func FindNumWordMFCC(words []Word) int {
	var result int
	for _, wd := range words {
		cols := wd.MFCCCols
		if cols > 0 {
			result = cols
			break
		}
	}
	return result
}
