package req

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func DecodeFile(path string) Request {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	resp := decode(content)
	return resp
}

func DecodeString(str string) Request {
	resp := decode([]byte(str))
	return resp
}

func decode(requestYaml []byte) Request {
	var resp Request
	err := yaml.Unmarshal(requestYaml, &resp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return resp
}
