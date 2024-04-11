package req

import (
	"fmt"
	"reflect"
)

func Validate(req Request) {
	var msgs []string
	checkRequired(req.Required, &msgs)
	checkTestament(&req.Testament)
	checkAudioData(&req.AudioData, &msgs)
	checkTextData(&req.TextData, &msgs)
	checkDetail(&req.Detail)
	checkTimestamps(&req.Timestamps, &msgs)
	checkAudioEncoding(&req.AudioEncoding, &msgs)
	checkTextEncoding(&req.TextEncoding, &msgs)
	checkOutputFormat(&req.OutputFormat, &msgs)
	checkCompare(req.Compare, &msgs)
	checkCompareSettings(req.Compare.CompareSettings, &msgs)
	fmt.Println("")
	for _, msg := range msgs {
		fmt.Println(msg)
	}
	fmt.Println("Testament", req.Testament)
	fmt.Println("Detail", req.Detail)
	fmt.Println("Timestamps", req.Timestamps)
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

func checkDetail(req *Detail) {
	if !req.Lines && !req.Words {
		req.Lines = true
	}
}

func checkTimestamps(req *Timestamps, msgs *[]string) {
	structVal := reflect.ValueOf(*req)
	count := checkForOne(structVal, msgs)
	if count == 0 {
		req.NoTimestamps = true
	}
}

func checkAudioEncoding(req *AudioEncoding, msgs *[]string) {
	structVal := reflect.ValueOf(*req)
	count := checkForOne(structVal, msgs)
	if count == 0 {
		req.NoEncoding = true
	}
}

func checkTextEncoding(req *TextEncoding, msgs *[]string) {
	structVal := reflect.ValueOf(*req)
	count := checkForOne(structVal, msgs)
	if count == 0 {
		req.NoEncoding = true
	}
}

func checkOutputFormat(req *OutputFormat, msgs *[]string) {
	structVal := reflect.ValueOf(*req)
	count := checkForOne(structVal, msgs)
	if count == 0 {
		req.JSON = true
	}
}

func checkCompare(req Compare, msgs *[]string) {
	if (req.Project1 != `` && req.Project2 == ``) ||
		(req.Project2 != `` && req.Project1 == ``) {
		*msgs = append(*msgs, `Compare must have two projects, not one.`)
	}
}

func checkCompareSettings(req CompareSettings, msgs *[]string) {
	checkForOne(reflect.ValueOf(req.DoubleQuotes), msgs)
	structVal := reflect.ValueOf(req.Apostrophe)
	checkForOne(structVal, msgs)
	checkForOne(reflect.ValueOf(req.Hyphen), msgs)
	checkForOne(reflect.ValueOf(req.DiacriticalMarks), msgs)
}

func checkForOne(structVal reflect.Value, msgs *[]string) int {
	var errorCount int
	var wasSet []string
	checkForOneRecursive(structVal, &wasSet)
	if len(wasSet) > 1 {
		for _, item := range wasSet {
			errorCount++
			msg := `Only 1 Data field can be set ` + item
			*msgs = append(*msgs, msg)
		}
	}
	return errorCount
}

func checkForOneRecursive(sVal reflect.Value, wasSet *[]string) {
	for i := 0; i < sVal.NumField(); i++ {
		field := sVal.Field(i)
		if field.Kind() == reflect.String {
			if field.String() != `` {
				*wasSet = append(*wasSet, sVal.Type().Name()+`.`+sVal.Type().Field(i).Name)
			}
		} else if field.Kind() == reflect.Bool {
			if field.Bool() {
				*wasSet = append(*wasSet, sVal.Type().Name()+`.`+sVal.Type().Field(i).Name)
			}
		} else if field.Kind() == reflect.Struct {
			checkForOneRecursive(field, wasSet)
		} else {
			msg := sVal.Type().Field(i).Name + ` has unexpected type ` + field.Type().Name()
			panic(msg)
		}
	}
}
