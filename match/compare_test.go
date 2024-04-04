package match

import (
	"context"
	"fmt"
	"testing"
)

func TestCompare(t *testing.T) {
	ctx := context.Background()
	compare := NewCompare(ctx, `ATIWBT_USXEDIT.db`, `ATIWBT_SCRIPT.db`)
	status := compare.Process()
	fmt.Println(status)
	if globalDiffCount != 2 {
		t.Error(`Expected count is 2, actual was`, globalDiffCount)
	}
}

/*if
func getCommandLine() (string, string) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: $HOME/Documents/go2/bin/compare  database1  database2")
		os.Exit(1)
	}
	return os.Args[1], os.Args[2]
}
*/
