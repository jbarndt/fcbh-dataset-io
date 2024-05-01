package output

import (
	"bytes"
	"dataset"
	log "dataset/logger"
	"encoding/csv"
	"encoding/json"
	"strconv"
)

func (o *Output) JSONStatus(status dataset.Status, debug bool) string {
	var result string
	if !debug {
		status.Trace = ``
	}
	bytes, err := json.Marshal(status)
	//err = errors.New(`wwww`) // for testing error path
	if err == nil {
		result = string(bytes)
	} else {
		status2 := log.Error(o.ctx, 500, err, `Error while creating error output`)
		result = status.String() + `, ` + status2.String()
	}
	return `[` + result + `]`
}

func (o *Output) CSVStatus(status dataset.Status, debug bool) string {
	var result string
	var buffer = bytes.NewBufferString("")
	writer := csv.NewWriter(buffer)
	_ = writer.Write([]string{`Name`, `Value`})
	_ = writer.Write([]string{`is_error`, strconv.FormatBool(status.IsErr)})
	_ = writer.Write([]string{`message`, status.Message})
	_ = writer.Write([]string{`status`, strconv.Itoa(status.Status)})
	_ = writer.Write([]string{`error`, status.Err})
	_ = writer.Write([]string{`request`, status.Request})
	if debug {
		_ = writer.Write([]string{`trace`, status.Trace})
	}
	writer.Flush()
	err := writer.Error()
	if err != nil {
		status2 := log.Error(o.ctx, 500, err, `Error while creating error output`)
		result = status.String() + `, ` + status2.String()
	}
	result = buffer.String()
	return result
}
