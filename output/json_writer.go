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
	_, _ = writer.WriteString(`[`)
	numStructs := len(structs)
	for line, scr := range structs {
		str := reflect.ValueOf(scr)
		_, _ = writer.WriteString(`{ `)
		for col, mt := range meta {
			data := str.Field(mt.Index)
			if data.Kind() == reflect.Slice {
				for i := 0; i < data.Len(); i++ {
					item := data.Index(i)
					if item.Kind() == reflect.Slice {
						if i > 0 {
							_, _ = writer.WriteString(`{ `)
						}
						for j := 0; j < item.Len(); j++ {
							write(writer, j, names[mt.CSVPos+j], ToString(item.Index(j)), mt)
						}
						if i < data.Len()-1 {
							_, _ = writer.WriteString(" },\n")
						}
					} else {
						write(writer, col, names[mt.CSVPos+i], ToString(item), mt)
					}
				}
			} else {
				write(writer, col, names[mt.CSVPos], ToString(data), mt)
			}
		}
		if line < numStructs-1 {
			_, _ = writer.WriteString(" },\n")
		} else {
			_, _ = writer.WriteString(` }`)
		}
	}
	_, _ = writer.WriteString("]\n")
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
	_ = file.Close()
	return file.Name()
}

func write(writer *bufio.Writer, col int, name string, value string, meta Meta) {
	if col > 0 {
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
