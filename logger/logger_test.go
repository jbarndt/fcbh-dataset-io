package logger

import (
	"context"
	"dataset"
	"testing"
)

func TestLogger_Main(t *testing.T) {
	Warn(context.Background(), 0, "Sample Error")
	Error(context.Background(), "Error Message")
}

func TestPanic(t *testing.T) {
	Panic(context.Background(), "Panic Message")
}

func TestFatal(t *testing.T) {
	ctx := context.Background()
	Fatal(ctx, "Error Message")
}

func TestError(t *testing.T) {
	req := dataset.RequestType{AudioSource: `ATWWBT2DAN`, TextSource: `ATIWBTN_T-USX`,
		Testament: dataset.NT}
	ctx := context.WithValue(context.Background(), "request", req)
	Error(ctx, "Error Message")
}
