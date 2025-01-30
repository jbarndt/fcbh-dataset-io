package stdio_exec

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestStdioExec(t *testing.T) {
	ctx := context.Background()
	uromanPath := os.Getenv(`FCBH_UROMAN_EXE`)
	stdio, status := NewStdioExec(ctx, uromanPath)
	result, status2 := stdio.Process("abc")
	fmt.Println("result:", result, status, status2)
}
