package encode

import (
	"context"
	"testing"
)

func TestAeneasExperiment(t *testing.T) {
	ctx := context.Background()
	aen := NewAeneasExperiment(ctx, `APFCMUN2DA`, `apf`)
	aen.Process()
}
