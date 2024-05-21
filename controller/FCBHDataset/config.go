package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Host       string
	BBKey      string
	AWSProfile string
}

func GetConfig() Config {
	var cfg Config
	homeDir, err := os.UserHomeDir()
	Catch(err)
	var file *os.File
	cfgPath := filepath.Join(homeDir, "FCBHDataset.yaml")
	_, err = os.Stat(cfgPath)
	// Read Config
	if err == nil || !os.IsNotExist(err) {
		file, err = os.Open(cfgPath)
		Catch(err)
		decoder := yaml.NewDecoder(file)
		decoder.KnownFields(true)
		err = decoder.Decode(&cfg)
		Catch(err)
		err = file.Close()
		Catch(err)
	}
	isChanged := false
	if cfg.Host == `` {
		cfg.Host = Prompt(`Host Address`)
		isChanged = true
	}
	if cfg.BBKey == `` {
		cfg.BBKey = Prompt(`Bible Brain Key`)
		isChanged = true
	}
	if cfg.AWSProfile == `` {
		cfg.AWSProfile = Prompt(`AWS Profile`)
		isChanged = true
	}
	// Save Config
	if isChanged {
		bytes, err := yaml.Marshal(&cfg)
		Catch(err)
		file, err = os.OpenFile(cfgPath, os.O_WRONLY|os.O_CREATE, 0666)
		Catch(err)
		_, err = file.Write(bytes)
		Catch(err)
		_ = file.Close()

	}
	return cfg
}

func Prompt(prompt string) string {
	fmt.Print(`Enter `, prompt, ` : `)
	var answer string
	count, err := fmt.Scanln(&answer)
	if count > 0 {
		Catch(err)
	}
	return answer
}
