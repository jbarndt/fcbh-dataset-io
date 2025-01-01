package tests

import (
	"dataset/controller"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRandomDirect(t *testing.T) {
	filename := "N2MKDMBS_proof.yaml"
	filePath := filepath.Join(os.Getenv("HOME"), filename)
	request, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	database, status := controller.CLIProcessEntry(request)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println("Test output", database)
}
