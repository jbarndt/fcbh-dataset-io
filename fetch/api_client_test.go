package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestAPIClient(t *testing.T) {
	ctx := context.Background()
	url := "https://4.dbt.io/api/languages/apf?v=4"
	content, status := HttpGet(ctx, url, "test")
	if status != nil {
		t.Fatal(status)
	}
	//fmt.Println(string(content))
	var response = make(map[string]map[string]any)
	err := json.Unmarshal(content, &response)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(response)
}
