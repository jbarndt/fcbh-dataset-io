package diff

import (
	"dataset/db"
	"dataset/decode_yaml/request"
	log "dataset/logger"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func (c *Compare) CompareScriptLines() (string, *log.Status) {
	var filename string
	var status *log.Status
	filename = c.writer.WriteHeading(c.baseDataset)
	// Build a map of the scripts to be compared
	var scripts []db.Script
	scripts, status = c.database.SelectScripts()
	if status != nil {
		return filename, status
	}
	var compareMap = make(map[string]db.Script)
	for _, comp := range scripts {
		compareMap[comp.ScriptNum] = comp
	}
	// Select and increment over the base scripts
	scripts, status = c.baseDb.SelectScripts()
	if status != nil {
		return filename, status
	}
	for _, base := range scripts {
		baseText := c.cleanup(base.ScriptText)
		if c.baseIdent.TextSource == request.TextScript {
			baseText = c.verseRm.ReplaceAllString(baseText, ``)
		}
		comp, ok := compareMap[base.ScriptNum]
		if ok {
			compText := c.cleanup(comp.ScriptText)
			if c.compIdent.TextSource == request.TextScript {
				compText = c.verseRm.ReplaceAllString(compText, ``)
			}
			c.scriptLineDiff(base.BookId, base.ChapterNum, base.ScriptNum, baseText, compText)
			delete(compareMap, base.ScriptNum)
		} else {
			c.scriptLineDiff(base.BookId, base.ChapterNum, base.ScriptNum, baseText, "")
		}
	}
	// Report any compare scripts that had no entry in base
	for _, comp := range compareMap {
		compText := c.cleanup(comp.ScriptText)
		if c.compIdent.TextSource == request.TextScript {
			compText = c.verseRm.ReplaceAllString(compText, ``)
		}
		c.scriptLineDiff(comp.BookId, comp.ChapterNum, comp.ScriptNum, "", compText)
	}
	filenameMap, status := c.generateBookChapterFilenameMap()
	c.baseDb.Close()
	c.writer.WriteEnd(filenameMap, c.insertSum, c.deleteSum, c.diffCount)
	return filename, status
}

func (c *Compare) scriptLineDiff(bookId string, chapter int, line string, baseText string, text string) {
	diffMatch := diffmatchpatch.New()
	diffs := diffMatch.DiffMain(baseText, text, false)
	if !c.isMatch(diffs) {
		inserts, deletes := c.measure(diffs)
		c.insertSum += inserts
		c.deleteSum += deletes
		avgLen := float64(len(baseText)+len(text)) / 2.0
		errPct := float64((inserts+deletes)*100) / avgLen
		c.writer.WriteScriptLineDiff(bookId, chapter, line, inserts, deletes, errPct, diffMatch.DiffPrettyHtml(diffs))
		c.diffCount++
	}
}
