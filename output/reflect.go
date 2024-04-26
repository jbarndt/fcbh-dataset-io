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

func FindNumScriptMFCC(structs []Script) int {
	var result int
	for _, str := range structs {
		if str.MFCCCols > 0 {
			result = str.MFCCCols
			break
		}
	}
	return result
}

func FindNumWordMFCC(structs []Word) int {
	var result int
	for _, str := range structs {
		if str.MFCCCols > 0 {
			result = str.MFCCCols
			break
		}
	}
	return result
}

func SetNumMFCC(meta *[]Meta, numMFCC int) {
	for i, mt := range *meta {
		if mt.Name == `MFCC` {
			(*meta)[i].Cols = numMFCC
		}
	}
}

func ConvertScriptsAny(structs []Script) []any {
	var results = make([]any, 0, len(structs))
	for _, str := range structs {
		results = append(results, str)
	}
	return results
}

func ConvertWordsAny(structs []Word) []any {
	var results = make([]any, 0, len(structs))
	for _, str := range structs {
		results = append(results, str)
	}
	return results
}

func FindActiveCols(structs []any, meta []Meta) []Meta {
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
