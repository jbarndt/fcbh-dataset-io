package request

import (
	"dataset"
	"reflect"
	"strings"
)

func (r *RequestDecoder) Validate(req *Request) dataset.Status {
	var msgs []string
	checkRequired(req.Required, &msgs)
	checkTestament(&req.Testament)
	checkAudioData(&req.AudioData, `AudioData`, &msgs)
	checkTextData(&req.TextData, `TextData`, &msgs)
	checkDetail(&req.Detail)
	checkTimestamps(&req.Timestamps, `Timestamps`, &msgs)
	checkAudioEncoding(&req.AudioEncoding, `AudioEncoding`, &msgs)
	checkTextEncoding(&req.TextEncoding, `TextEncoding`, &msgs)
	checkOutputFormat(&req.OutputFormat, `OutputFormat`, &msgs)
	//checkCompare(req.Compare, &msgs)
	checkForOne(reflect.ValueOf(req.Compare.CompareSettings.DoubleQuotes), `DoubleQuotes`, &msgs)
	checkForOne(reflect.ValueOf(req.Compare.CompareSettings.Apostrophe), `Apostrophe`, &msgs)
	checkForOne(reflect.ValueOf(req.Compare.CompareSettings.Hyphen), `Hyphen`, &msgs)
	checkForOne(reflect.ValueOf(req.Compare.CompareSettings.DiacriticalMarks), `DiscriticalMarks`, &msgs)
	//fmt.Println("")
	var status dataset.Status
	if len(msgs) > 0 {
		status.Status = 400
		status.IsErr = true
		status.Message = strings.Join(msgs, "\n")
		//status.Request =
	}
	return status
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
	if !req.OT && !req.NT && len(req.NTBooks) == 0 && len(req.OTBooks) == 0 {
		req.NT = true
	}
}

// checkAudioData Is checking that no more than one item is selected.
// if none are selected, it will set the default: NoAudio
func checkAudioData(req *AudioData, fieldName string, msgs *[]string) {
	count := checkForOne(reflect.ValueOf(*req), fieldName, msgs)
	if count == 0 {
		req.NoAudio = true
	}
}

// checkTextData Is checking that no more than one item is selected.
// if none are selected, it will set the default: NoAudio
func checkTextData(req *TextData, fieldName string, msgs *[]string) {
	count := checkForOne(reflect.ValueOf(*req), fieldName, msgs)
	if count == 0 {
		req.NoText = true
	}
}

func checkDetail(req *Detail) {
	if !req.Lines && !req.Words {
		req.Lines = true
	}
}

func checkTimestamps(req *Timestamps, fieldName string, msgs *[]string) {
	count := checkForOne(reflect.ValueOf(*req), fieldName, msgs)
	if count == 0 {
		req.NoTimestamps = true
	}
}

func checkAudioEncoding(req *AudioEncoding, fieldName string, msgs *[]string) {
	count := checkForOne(reflect.ValueOf(*req), fieldName, msgs)
	if count == 0 {
		req.NoEncoding = true
	}
}

func checkTextEncoding(req *TextEncoding, fieldName string, msgs *[]string) {
	count := checkForOne(reflect.ValueOf(*req), fieldName, msgs)
	if count == 0 {
		req.NoEncoding = true
	}
}

func checkOutputFormat(req *OutputFormat, fieldName string, msgs *[]string) {
	count := checkForOne(reflect.ValueOf(*req), fieldName, msgs)
	if count == 0 {
		req.JSON = true
	}
}

//func checkCompare(req Compare, msgs *[]string) {
//	if (req.Project1 != `` && req.Project2 == ``) ||
//		(req.Project2 != `` && req.Project1 == ``) {
//		*msgs = append(*msgs, `Compare must have two projects, not one.`)
//	}
//}

func checkForOne(structVal reflect.Value, fieldName string, msgs *[]string) int {
	var errorCount int
	var wasSet []string
	checkForOneRecursive(structVal, &wasSet)
	if len(wasSet) > 1 {
		errorCount += len(wasSet)
		msg := `Only 1 field can be set on ` + fieldName + `: ` + strings.Join(wasSet, `,`)
		*msgs = append(*msgs, msg)
	}
	return errorCount
}

func checkForOneRecursive(sVal reflect.Value, wasSet *[]string) {
	for i := 0; i < sVal.NumField(); i++ {
		field := sVal.Field(i)
		if field.Kind() == reflect.String {
			if field.String() != `` {
				*wasSet = append(*wasSet, sVal.Type().Field(i).Name)
			}
		} else if field.Kind() == reflect.Bool {
			if field.Bool() {
				*wasSet = append(*wasSet, sVal.Type().Field(i).Name)
			}
		} else if field.Kind() == reflect.Struct {
			checkForOneRecursive(field, wasSet)
		} else {
			msg := sVal.Type().Field(i).Name + ` has unexpected type ` + field.Type().Name()
			panic(msg)
		}
	}
}
