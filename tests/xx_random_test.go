package tests

import (
	"dataset/controller"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRandomDirect(t *testing.T) {
	//filename := "N2MKDMBS_proof.yaml"
	//filename := "N2HOYWFW_proof.yaml"
	//filename := "N2ENGWEB_proof.yaml"
	filename := "N2CUL_MNT_rpt.yaml"
	filePath := filepath.Join(os.Getenv("HOME"), filename)
	request, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	database, status := controller.CLIProcessEntry(request)
	if status != nil {
		t.Fatal(status)
	}
	fmt.Println("Test output", database)
}
