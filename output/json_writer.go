package output

import (
	"bufio"
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
	for line, scr := range structs {
		str := reflect.ValueOf(scr)
		for col, mt := range meta {
			data := str.Field(mt.Index)
			if data.Kind() == reflect.Slice || data.Kind() == reflect.Array {
				for i := 0; i < data.Len(); i++ {
					item := data.Index(i)
					if item.Kind() == reflect.Slice || item.Kind() == reflect.Array {
						for j := 0; j < item.Len(); j++ {
							write(writer, line+i, j, names[col+j], ToString(item.Index(j)), mt)
						}
					} else {
						write(writer, line+i, col, names[col], ToString(item), mt)
					}
				}
			} else {
				write(writer, line, col, names[col], ToString(data), mt)
			}
		}
	}
	_, _ = writer.WriteString(" }]\n")
	_ = file.Close()
	return file.Name()
}

func write(writer *bufio.Writer, line int, col int, name string, value string, meta Meta) {
	if line == 0 && col == 0 {
		_, _ = writer.WriteString(`[{ "`)
	} else if line > 0 && col == 0 {
		_, _ = writer.WriteString(" },\n{ ")
	} else {
		_, _ = writer.WriteString(`, "`)
	}
	_, _ = writer.WriteString(name)
	_, _ = writer.WriteString(`": `)
	if meta.Dtype == `string` {
		_, _ = writer.WriteString(`"`)
		_, _ = writer.WriteString(value)
		_, _ = writer.WriteString(`"`)
	} else {
		_, _ = writer.WriteString(value)
	}
}
