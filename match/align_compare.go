package match

import (
	"dataset"
	"dataset/db"
	"dataset/generic"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"strings"
)

func (a *AlignErrorCalc) compareLines2ASR(lines []generic.AlignLine, asrConn db.DBAdapter) ([]generic.AlignLine, dataset.Status) {
	var result []generic.AlignLine
	var status dataset.Status
	for _, line := range lines {
		var silencePos = a.FindSilencePos(line.Chars)
		if len(silencePos) == 0 {
			result = append(result, line)
		} else {
			var newLine generic.AlignLine
			lineRef := line.Chars[0].LineRef
			var asrText string
			asrText, status = asrConn.SelectLine(lineRef)
			if status.IsErr {
				return result, status
			}
			alignNorm, alignUroman := a.GetOriginalText(line.Chars)
			fmt.Println(len(alignUroman))
			cDiffs := a.DiffMatchPatch(lineRef, alignNorm, asrText)
			if len(cDiffs) == 0 {
				continue
			}
			for _, silPos := range silencePos {
				diffPos := a.FindPositionInDiff(cDiffs, silPos)
				//for _, ch := range line.Chars {
				for i := 0; i < silPos; i++ {
					newLine.Chars = append(newLine.Chars, line.Chars[i])
				}
				curr := line.Chars[silPos]
				fmt.Println(silPos, string(alignNorm[silPos]), string(cDiffs[diffPos].Char))
				// Add every insert character after this
				for i := diffPos + 1; i < len(cDiffs); i++ {
					if cDiffs[i].Type == diffmatchpatch.DiffInsert {
						fmt.Println("add char ASR char", string(cDiffs[i].Char))
						var newChar generic.AlignChar
						newChar.AudioFile = curr.AudioFile
						newChar.LineId = curr.LineId
						newChar.LineRef = curr.LineRef
						//newChar.WordId = ch.WordId // might not be correct
						newChar.CharNorm = cDiffs[i].Char
						//newChar.CharUroman =
						//newChar.BeginTS = beginTS
						//newChar.EndTS = endTS // It a number of chars are found they have the same TS
						//newChar.FAScore = 1.0
						newChar.IsASR = true
						newLine.Chars = append(newLine.Chars, newChar)
					} else {
						break
					}
				}
			}
			for i := silencePos[len(silencePos)-1]; i < len(newLine.Chars); i++ {
				newLine.Chars = append(newLine.Chars, line.Chars[i])
			}
			result = append(result, newLine)
		}
	}
	return result, status
}

func (a *AlignErrorCalc) FindSilencePos(chars []generic.AlignChar) []int {
	var silencePos []int
	var pos = -1
	for i, char := range chars {
		if i > 0 && char.CharSeq == 0 {
			pos += 2
		} else {
			pos += 1
		}
		if char.SilenceLong > 0 {
			silencePos = append(silencePos, pos)
		}
	}
	return silencePos
}

func (a *AlignErrorCalc) GetOriginalText(chars []generic.AlignChar) (string, string) {
	var alNorm []rune
	var alUroman []rune
	for i, char := range chars {
		if i > 0 && char.CharSeq == 0 {
			alNorm = append(alNorm, ' ')
			alUroman = append(alUroman, ' ')
		}
		alNorm = append(alNorm, char.CharNorm)
		alUroman = append(alUroman, char.CharUroman)
	}
	alignNorm := strings.ToLower(string(alNorm))
	alignUroman := string(alUroman)
	return alignNorm, alignUroman
}

type CDiff struct {
	Type diffmatchpatch.Operation
	Char rune
}

func (a *AlignErrorCalc) DiffMatchPatch(lineRef string, text string, asrText string) []CDiff {
	var result []CDiff
	diffMatch := diffmatchpatch.New()
	text = strings.TrimSpace(text)
	asrText = strings.TrimSpace(asrText)
	diffs := diffMatch.DiffMain(text, asrText, false)
	diffs = diffMatch.DiffCleanupSemantic(diffs)
	if len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual {
		return result
	}
	fmt.Println(lineRef, asrText)
	fmt.Println(lineRef, text)
	fmt.Println(lineRef, diffMatch.DiffPrettyText(diffs))
	fmt.Println(diffs)
	for _, df := range diffs {
		for _, ch := range df.Text {
			var cDiff CDiff
			cDiff.Type = df.Type
			cDiff.Char = ch
			result = append(result, cDiff)
		}
	}
	return result
}

func (a *AlignErrorCalc) FindPositionInDiff(cDiffs []CDiff, charPos int) int {
	var diffPos = -1
	for i, ch := range cDiffs {
		if ch.Type != diffmatchpatch.DiffInsert {
			diffPos++
			if diffPos >= charPos {
				return i
			}
		}
	}
	return len(cDiffs)
}
