package output

import (
	"reflect"
	"strings"
)

func (o *Output) ReflectStruct(aStruct any) []Meta {
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

func (o *Output) FindNumScriptMFCC(structs []Script) int {
	var result int
	for _, str := range structs {
		if str.MFCCCols > 0 {
			result = str.MFCCCols
			break
		}
	}
	return result
}

func (o *Output) FindNumWordMFCC(structs []Word) int {
	var result int
	for _, str := range structs {
		if str.MFCCCols > 0 {
			result = str.MFCCCols
			break
		}
	}
	return result
}

func (o *Output) SetNumMFCC(meta *[]Meta, numMFCC int) {
	for i, mt := range *meta {
		if mt.Name == `MFCC` {
			(*meta)[i].Cols = numMFCC
		}
	}
}

func (o *Output) FindNumWordEnc(structs []Word, meta *[]Meta) {
	var length int
	for _, str := range structs {
		if len(str.WordEnc) > 0 {
			length = len(str.WordEnc)
			break
		}
	}
	for i, mt := range *meta {
		if mt.Name == `WordEnc` {
			(*meta)[i].Cols = length
		}
	}
}

func (o *Output) SetCSVPos(meta *[]Meta) {
	var pos = 0
	for i, str := range *meta {
		(*meta)[i].CSVPos = pos
		pos += str.Cols
	}
}

func (o *Output) ConvertScriptsAny(structs []Script) []any {
	var results = make([]any, 0, len(structs))
	for _, str := range structs {
		results = append(results, str)
	}
	return results
}

func (o *Output) ConvertWordsAny(structs []Word) []any {
	var results = make([]any, 0, len(structs))
	for _, str := range structs {
		results = append(results, str)
	}
	return results
}

func (o *Output) FindActiveCols(structs []any, meta []Meta) []Meta {
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
