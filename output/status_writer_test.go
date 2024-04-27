package output

import (
	"context"
	"dataset"
	log "dataset/logger"
	"dataset/request"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestDefaultStatus(t *testing.T) {
	ctx := context.Background()
	status := log.Error(ctx, 400, errors.New(`Huh`), `My message`)
	str := status.String()
	if len(str) != 76 {
		t.Error(`Status should be len 76`, len(str))
	}
	var response any
	err := json.Unmarshal([]byte(str), &response)
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(response)
}

func TestJSONStatus(t *testing.T) {
	status, ctx := prepareError(t)
	result := JSONStatus(ctx, status, true)
	fmt.Println(result)
	if len(result) != 360 {
		t.Error(`Result should be len 360`, len(result))
	}
	var response any
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Error(err)
	}
}

func TestCSVStatus(t *testing.T) {
	status, ctx := prepareError(t)
	result := CSVStatus(ctx, status, true)
	if len(result) != 342 {
		t.Error(`Result should be len 342`, len(result))
	}
	//fmt.Println(result)
}

func prepareError(t *testing.T) (dataset.Status, context.Context) {
	var req request.Request
	req.Required.RequestName = `Test1`
	req.Required.BibleId = `ENGWEB`
	req.Required.IsNew = true
	req.Required.LanguageISO = `eng`
	req.Testament.NT = true
	ctx := context.Background()
	reqDecoder := request.NewRequestDecoder(ctx)
	yaml, status := reqDecoder.Encode(req)
	if status.IsErr {
		t.Error(status.Message)
	}
	ctx = context.WithValue(context.Background(), `request`, yaml)
	err := errors.New("test err")
	status = log.Error(ctx, 400, err, "my error message")
	return status, ctx
}
