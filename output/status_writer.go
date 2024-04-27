package output

import (
	"context"
	"dataset"
	log "dataset/logger"
	"encoding/json"
)

func JSONStatus(ctx context.Context, status dataset.Status, debug bool) string {
	var result string
	if !debug {
		status.Trace = ``
	}
	bytes, err := json.Marshal(status)
	//err = errors.New(`wwww`) // for testing error path
	if err == nil {
		result = string(bytes)
	} else {
		status2 := log.Error(ctx, 500, err, `Error while creating error output`)
		result = status.String() + `, ` + status2.String()
	}
	return `[` + result + `]`
}

func CSVStatus(ctx context.Context, status dataset.Status, debug bool) string {
	var result string
	return result
}
