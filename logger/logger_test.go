package logger

import (
	"context"
	"testing"
)

func TestLogger_Main(t *testing.T) {
	Warn(context.Background(), 0, "Sample Error")
}
