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
func checkAudioData(req *AudioData, msgs *[]string) int {
	var errorCount int
	var wasSet []string
	sVal := reflect.ValueOf(*req)
	for i := 0; i < sVal.NumField(); i++ {
		field := sVal.Field(i)
		fType := field.Type().Name()
		if fType == `string` {
			if field.String() != `` {
				sTyp := sVal.Type()
				wasSet = append(wasSet, sTyp.Name()+`.`+sTyp.Field(i).Name)
			}
		} else if fType == `bool` {
			if field.Bool() {
				sTyp := sVal.Type()
				wasSet = append(wasSet, sTyp.Name()+`.`+sTyp.Field(i).Name)
			}
		} else if fType == `BibleBrainAudio` {
			bbVal := reflect.ValueOf(req.BibleBrain)
			for j := 0; j < bbVal.NumField(); j++ {
				bbField := bbVal.Field(j)
				//bbType := bbField.Type().Name()
				if bbField.Bool() {
					bbType := bbVal.Type()
					wasSet = append(wasSet, bbType.Name()+`.`+bbType.Field(j).Name)
				}
			}
		} else {
			msg := sVal.Type().Field(i).Name + ` has unexpected type ` + fType
			panic(msg)
		}
	}
	if len(wasSet) > 0 {
		for _, item := range wasSet {
			errorCount++
			msg := `Only 1 Data field can be set ` + item
			*msgs = append(*msgs, msg)
		}
	} else if len(wasSet) == 0 {
		req.NoAudio = true
	}
	return errorCount
}
