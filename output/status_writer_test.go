package output

import (
	"context"
	log "dataset/logger"
	"dataset/request"
	"errors"
	"testing"
)

func TestDefaultStatus(t *testing.T) {
	ctx := context.Background()
	status := log.Error(ctx, 400, errors.New(`Huh`), `My message`)
	str := status.String()
	if len(str) != 76 {
		t.Error(`Status should be len 76`, len(str))
	}
}

func TestWriteJSONStatus(t *testing.T) {
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
	result := JSONStatus(ctx, status, true)
	if len(result) != 318 {
		t.Error(`Result should be len 318`, len(result))
	}
}
