package req

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"reflect"
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

func Validate(req Request) {
	var msgs []string
	checkRequired(req.Required, &msgs)
	checkTestament(&req.Testament)
	checkAudioData(&req.AudioData, &msgs)
	//checkTextData(&req.TextData, &msgs)
	for _, msg := range msgs {
		fmt.Println(msg)
	}
	fmt.Println(req.Testament)
}

func checkRequired(req Required, msgs *[]string) {
	sVal := reflect.ValueOf(req)
	for i := 0; i < sVal.NumField(); i++ {
		field := sVal.Field(i)
		if field.String() == `` {
			sTyp := sVal.Type()
			name := sTyp.Name() + `.` + sTyp.Field(i).Name
			msg := `Required field ` + name + ` is missing.`
			*msgs = append(*msgs, msg)
		}
	}
}

func checkTestament(req *Testament) {
	if !req.OT && !req.NT {
		req.NT = true
	}
}

// checkAudioData Is checking that no more than one item is selected.
// if none are selected, it will set the default: NoAudio
func checkAudioData(req *AudioData, msgs *[]string) {
	sVal := reflect.ValueOf(*req)
	count := checkForOne(sVal, msgs)
	if count == 0 {
		req.NoAudio = true
	}
}

// checkTextData Is checking that no more than one item is selected.
// if none are selected, it will set the default: NoAudio
func checkTextData(req *TextData, msgs *[]string) {
	sVal := reflect.ValueOf(*req)
	count := checkForOne(sVal, msgs)
	if count == 0 {
		req.NoText = true
	}
}

func checkForOne(sVal reflect.Value, msgs *[]string) int {
	var errorCount int
	var wasSet []string
	for i := 0; i < sVal.NumField(); i++ {
		field := sVal.Field(i)
		if field.Kind() == reflect.String {
			if field.String() != `` {
				wasSet = append(wasSet, sVal.Type().Name()+`.`+sVal.Type().Field(i).Name)
			}
		} else if field.Kind() == reflect.Bool {
			if field.Bool() {
				wasSet = append(wasSet, sVal.Type().Name()+`.`+sVal.Type().Field(i).Name)
			}
		} else if field.Kind() == reflect.Struct {
			for j := 0; j < field.NumField(); j++ {
				bbField := field.Field(j)
				if bbField.Bool() {
					wasSet = append(wasSet, field.Type().Name()+`.`+field.Type().Field(j).Name)
				}
			}
		} else {
			msg := sVal.Type().Field(i).Name + ` has unexpected type ` + field.Type().Name()
			panic(msg)
		}
	}
	if len(wasSet) > 0 {
		for _, item := range wasSet {
			errorCount++
			msg := `Only 1 Data field can be set ` + item
			*msgs = append(*msgs, msg)
		}
	}
	return errorCount
}
