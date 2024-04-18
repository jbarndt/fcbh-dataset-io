package fetch

import (
	"context"
	"fmt"
	"testing"
)

func TestAPIDownloadList(t *testing.T) {
	ctx := context.Background()
	bibleId := ``
	client := NewAPIDBPClient(ctx, bibleId)
	list, status := client.DownloadList()
	if status.IsErr {
		t.Error(status.Message)
	}
	for key, value := range list {
		fmt.Println(key, value)
	}
	fmt.Println(len(list))
}
