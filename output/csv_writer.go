package output

import (
	"encoding/csv"
	"os"
	"reflect"
	"strconv"
)

func WriteScriptCSV(scripts []Script, meta []Meta) string {
	var results = make([]any, 0, len(scripts))
	for _, script := range scripts {
		results = append(results, script)
	}
	return WriteCSV(results, meta)
}

func WriteCSV(structs []any, meta []Meta) string {
	file, err := os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), "csv")
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(file)
	var header []string
	for _, mt := range meta {
		if mt.Cols == 1 {
			header = append(header, mt.Tag)
		} else {
			for i := 0; i < mt.Cols; i++ {
				header = append(header, mt.Tag+strconv.Itoa(i))
			}
		}
	}
	_ = writer.Write(header)
	for _, scr := range structs {
		str := reflect.ValueOf(scr)
		var line []string
		for col, mt := range meta {
			data := str.Field(mt.Index)
			if data.Kind() == reflect.Slice || data.Kind() == reflect.Array {
				for i := 0; i < data.Len(); i++ {
					item := data.Index(i)
					if item.Kind() == reflect.Slice || item.Kind() == reflect.Array {
						if i > 0 {
							line = make([]string, col)
						}
						for j := 0; j < item.Len(); j++ {
							line = append(line, ToString(item.Index(j)))
						}
					} else {
						line = append(line, ToString(item))
					}
					_ = writer.Write(line)
				}
			} else {
				line = append(line, ToString(data))
			}
		}
		_ = writer.Write(line)
	}
	writer.Flush()
	err = writer.Error()
	if err != nil {
		panic(err)
	}
	_ = file.Close()
	return file.Name()
}

// ToString converts scalar values to string.  It does not convert the following kind.
// reflect.Uintptr, reflect.Complex64, reflect.Complex128, reflect.Struct, reflect.Array,
// reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface, reflect.Invalid
func ToString(value reflect.Value) string {
	var result string
	switch value.Kind() {
	case reflect.Bool:
		result = strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		result = strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.Float64:
		result = strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.String:
		result = value.String()
	default:
		panic(`output.ToString() cannot convert value of type` + value.Type().String())
	}
	return result
}
