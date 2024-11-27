package mms

import (
	"context"
	"fmt"
	"testing"
)

func TestURoman(t *testing.T) {
	var input []string
	ctx := context.Background()
	input = append(input, "Игорь Стравинский")
	input = append(input, "Игорь")
	input = append(input, "Ντέιβις Καπ")
	results, status := URoman(ctx, "rus", input)
	if status.IsErr {
		t.Fatal(status.String())
	}
	for _, line := range results {
		fmt.Println(line)
	}
}
