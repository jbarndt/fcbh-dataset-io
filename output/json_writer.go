package output

import (
	"bufio"
	"encoding/json"
	"os"
	"reflect"
	"strconv"
)

func WriteJSON(structs []any, meta []Meta) string {
	file, err := os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), "json")
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)
	var names []string
	for _, mt := range meta {
		if mt.Cols == 1 {
			names = append(names, mt.Tag)
		} else {
			for i := 0; i < mt.Cols; i++ {
				names = append(names, mt.Tag+strconv.Itoa(i))
			}
		}
	}
	var resp = make([]map[string]any, 0, len(structs))
	for _, scr := range structs {
		str := reflect.ValueOf(scr)
		var rec = make(map[string]any)
		for _, mt := range meta {
			data := str.Field(mt.Index)
			if data.Kind() == reflect.Slice {
				for i := 0; i < data.Len(); i++ {
					item := data.Index(i)
					if item.Kind() == reflect.Slice {
						for j := 0; j < item.Len(); j++ {
							rec[names[mt.CSVPos+j]] = ToValue(item.Index(j))
						}
						if i < data.Len()-1 {
							resp = append(resp, rec)
							rec = make(map[string]any)
						}
					} else {
						rec[names[mt.CSVPos+i]] = ToValue(item)
					}
				}
			} else {
				rec[names[mt.CSVPos]] = ToValue(data)
			}
		}
		resp = append(resp, rec)
		rec = make(map[string]any)
	}
	var encoder = json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(resp)
	if err != nil {
		panic(err)
	}
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
	_ = file.Close()
	return file.Name()
}

func ToValue(value reflect.Value) any {
	var result any
	switch value.Kind() {
	case reflect.Bool:
		result = value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = value.Uint()
	case reflect.Float32:
		result = value.Float()
	case reflect.Float64:
		result = value.Float()
	case reflect.String:
		result = value.String()
	default:
		panic(`output.ToValue() cannot convert value of type` + value.Type().String())
	}
	return result
}
