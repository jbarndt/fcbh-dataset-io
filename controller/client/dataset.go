package main

import (
	"dataset/controller"
	"fmt"
	"os"
)

func main() {
	// Should check args, and expect a filepath
	var yamlPath = `controller/client/request_test.yaml`
	var content, err = os.ReadFile(yamlPath)
	if err != nil {
		fmt.Println(err)
	} else {
		var control = controller.NewController(content)
		control.Process()
	}
}
