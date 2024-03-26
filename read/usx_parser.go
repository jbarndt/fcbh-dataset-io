package read

import (
	"dataset_io"
	"dataset_io/db"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Stack []string

var hasStyle = map[string]bool{
	`book`: true, `para`: true, `char`: true, `cell`: true, `ms`: true, `note`: true, `sidebar`: true, `figure`: true}

var numericPattern = regexp.MustCompile(`^\d+`)

// MOVE TO TEST FILE
//
//	func main() {
//		if len(os.Args) < 2 {
//			fmt.Println("Usage:  $HOME/Documents/go2/bin/usx_parser  bibleId")
//			os.Exit(1)
//		}
//		var bibleId = os.Args[1]
//		dbPath := os.Getenv(`FCBH_DATASET_DB`)
//		var db = openDatabase(dbPath, bibleId+"_USXEDIT.db")
//		directory := filepath.Join(dbPath, `download`, bibleId)
func ReadUSXEdit(database db.DBAdapter, bibleId string, testament dataset_io.TestamentType) {
	directory := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId)
	dirs, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	var suffix string
	switch testament {
	case dataset_io.NT:
		suffix = `N_ET-usx`
	case dataset_io.OT:
		suffix = `O_ET-usx`
	case dataset_io.ONT:
		suffix = `-usx`
	default:
		log.Fatal("Error: Unknown testament type", testament, "in ReadUSXEdit")
	}
	for _, dir := range dirs {
		if strings.HasSuffix(dir.Name(), suffix) {
			subDir := filepath.Join(directory, dir.Name())
			files, err := os.ReadDir(subDir)
			if err != nil {
				log.Fatal(err)
			}
			for _, file := range files {
				filename := filepath.Join(subDir, file.Name())
				fmt.Println(filename)
				records := decode(filename) // Also edits out non-script elements
				titleDesc := extractTitles(records)
				records = addChapterHeading(records, titleDesc)
				records = correctScriptNum(records)
				database.InsertScripts(records)
			}
		}
	}
	count := database.SelectScalarInt(`SELECT count(*) FROM scripts`)
	fmt.Println("Total Records", count)
}

func decode(filename string) []db.InsertScriptRec {
	xmlFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer xmlFile.Close()
	var stack Stack
	var rec db.InsertScriptRec
	var records []db.InsertScriptRec
	var tagName string
	var bookId string
	var chapterNum = 1
	var verseNum int
	var verseStr string
	var usfmStyle string
	decoder := xml.NewDecoder(xmlFile)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			panic(err)
		}
		switch se := token.(type) {
		//[StartElement], [EndElement], [CharData], [Comment], [ProcInst], [Directive].
		case xml.StartElement:
			tagName = se.Name.Local
			if tagName == `book` {
				bookId = findAttr(se, `code`)
			} else if tagName == `chapter` {
				chapterNum = findIntAttr(se, `number`)
			} else if tagName == `verse` {
				verseNum = findIntAttr(se, `number`)
				verseStr = findAttr(se, `number`)
			}
			if hasStyle[tagName] {
				usfmStyle = tagName + `.` + findAttr(se, `style`)
				if include(usfmStyle) { // This if fileters out the non-script code
					stack = stack.Push(usfmStyle)
				} else {
					err := decoder.Skip()
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		case xml.CharData:
			text := string(se)
			if len(strings.TrimSpace(text)) > 0 {
				// This is needed because scripts use {n} as a verse number.
				if strings.Contains(text, "{") || strings.Contains(text, "}") {
					fmt.Println("Found { } changed to ( ) in:", text)
					text = strings.Replace(text, `{`, `(`, -1)
					text = strings.Replace(text, `}`, `)`, -1)
				}
				rec.ScriptText = append(rec.ScriptText, text)
			}
		case xml.EndElement:
			if hasStyle[se.Name.Local] {
				stack, usfmStyle = stack.Pop()
			}
		}
		if chapterNum != rec.ChapterNum || verseNum != rec.VerseNum || usfmStyle != rec.UsfmStyle {
			if bookId != `` && len(rec.ScriptText) > 0 {
				records = append(records, rec)
			}
			scriptNum, err := strconv.Atoi(rec.ScriptNum)
			if err != nil {
				log.Fatalln(err)
			}
			scriptNum++
			if chapterNum != rec.ChapterNum {
				scriptNum = 1
			}
			rec = db.InsertScriptRec{BookId: bookId, ChapterNum: chapterNum, ScriptNum: strconv.Itoa(scriptNum),
				VerseNum: verseNum, VerseStr: verseStr, UsfmStyle: usfmStyle}
		}
	}
	fmt.Println("num records", len(records))
	return records
}

type titleDesc struct {
	heading string
	title   []db.InsertScriptRec
	areDiff bool
}

func extractTitles(records []db.InsertScriptRec) titleDesc {
	var results titleDesc
	for _, rec := range records {
		if rec.UsfmStyle == `para.h` {
			results.heading = strings.Join(rec.ScriptText, ``)
		} else if strings.HasPrefix(rec.UsfmStyle, `para.mt`) {
			results.title = append(results.title, rec)
		}
	}
	if results.heading == `` && len(results.title) > 0 {
		results.heading = strings.Join(results.title[len(results.title)-1].ScriptText, ``)
	}
	return results
}

func addChapterHeading(records []db.InsertScriptRec, titles titleDesc) []db.InsertScriptRec {
	var results = make([]db.InsertScriptRec, 0, len(records))
	for _, rec := range titles.title {
		results = append(results, rec)
	}
	var lastChapter = -1
	for _, rec := range records {
		if rec.ChapterNum != lastChapter {
			lastChapter = rec.ChapterNum
			var rec2 = rec
			rec2.VerseNum = 0
			rec2.VerseStr = ``
			rec2.ScriptText = []string{titles.heading + " " + strconv.Itoa(rec.ChapterNum)}
			results = append(results, rec2)
		}
		if rec.UsfmStyle != `para.h` && !strings.HasPrefix(rec.UsfmStyle, `para.mt`) {
			results = append(results, rec)
		}
	}
	return results
}

func correctScriptNum(records []db.InsertScriptRec) []db.InsertScriptRec {
	var results = make([]db.InsertScriptRec, 0, len(records))
	var scriptNum = 0
	var lastChapter = 0
	for _, rec := range records {
		if rec.ChapterNum != lastChapter {
			lastChapter = rec.ChapterNum
			scriptNum = 0
		}
		scriptNum += 1
		rec.ScriptNum = strconv.Itoa(scriptNum)
		results = append(results, rec)
	}
	return results
}

func findAttr(se xml.StartElement, name string) string {
	for _, attr := range se.Attr {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ``
}

func findIntAttr(se xml.StartElement, name string) int {
	val := findAttr(se, name)
	if val == `` {
		return 0
	} else {
		num := numericPattern.FindString(val)
		result, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		return result
	}
}

func (s Stack) Push(v string) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, string) {
	l := len(s)
	if l < 1 {
		panic("You tried to pop an empty stack.")
	}
	return s[:l-1], s[l-1]
}
