package match

import (
	"context"
	"dataset"
	log "dataset/logger"
	"os"
	"strconv"
	"time"
)

type HTMLWriter struct {
	ctx         context.Context
	datasetName string
	out         *os.File
}

func NewHTMLWriter(ctx context.Context, datasetName string) (HTMLWriter, dataset.Status) {
	var h HTMLWriter
	h.ctx = ctx
	h.datasetName = datasetName
	var status dataset.Status
	var err error
	h.out, err = os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), datasetName+"_*.html")
	if err != nil {
		status = log.Error(ctx, 500, err, `Error creating output file for diff`)
	}
	return h, status
}

func (h *HTMLWriter) WriteHeading(baseDataset string) {
	head := `<DOCTYPE html>
<html>
 <head>
  <meta charset="utf-8">
  <title>File Difference</title>
  <style>
p { margin: 20px 40px; }
  </style>
 </head>
 <body>`
	_, _ = h.out.WriteString(head)
	_, _ = h.out.WriteString(`<h2 style="text-align:center">Compare `)
	_, _ = h.out.WriteString(baseDataset)
	_, _ = h.out.WriteString(` to `)
	_, _ = h.out.WriteString(h.datasetName)
	_, _ = h.out.WriteString("</h2>\n")
	_, _ = h.out.WriteString(`<h3 style="text-align:center">`)
	_, _ = h.out.WriteString(time.Now().Format(`Mon Jan 2 2006 03:04:05 pm MST`))
	_, _ = h.out.WriteString("</h3>\n")
	_, _ = h.out.WriteString(`<h3 style="text-align:center">RED characters are those in `)
	_, _ = h.out.WriteString(baseDataset)
	_, _ = h.out.WriteString(` only, while GREEN characters are in `)
	_, _ = h.out.WriteString(h.datasetName)
	_, _ = h.out.WriteString(" only</h3>\n")
}

func (h *HTMLWriter) WriteVerseDiff(verse pair, inserts int, deletes int, diffHtml string) {
	ref := verse.bookId + " " + strconv.Itoa(verse.chapter) + ":" + verse.num + ` `
	//fmt.Println(ref, diffMatch.DiffPrettyText(diffs))
	//fmt.Println("=============")
	_, _ = h.out.WriteString(`<p>`)
	_, _ = h.out.WriteString(ref)
	_, _ = h.out.WriteString(` +`)
	_, _ = h.out.WriteString(strconv.Itoa(inserts))
	_, _ = h.out.WriteString(` -`)
	_, _ = h.out.WriteString(strconv.Itoa(deletes))
	_, _ = h.out.WriteString(` `)
	//_, _ = h.out.WriteString(diffMatch.DiffPrettyHtml(diffs))
	_, _ = h.out.WriteString(diffHtml)
	_, _ = h.out.WriteString("</p>\n")
	//_, _ = output.WriteString(`<p>`)
	//_, _ = output.WriteString(fmt.Sprint(diffs))
	//_, _ = output.WriteString("</p>\n")
}

func (h *HTMLWriter) WriteChapterDiff(bookId string, chapter int, inserts int, deletes int, diffHtml string) {
	_, _ = h.out.WriteString(`<h3 style="padding-left:50px;">`)
	_, _ = h.out.WriteString(bookId)
	_, _ = h.out.WriteString(" ")
	_, _ = h.out.WriteString(strconv.Itoa(chapter))
	_, _ = h.out.WriteString(", Inserts: ")
	_, _ = h.out.WriteString(strconv.Itoa(inserts))
	_, _ = h.out.WriteString(", Deletes: ")
	_, _ = h.out.WriteString(strconv.Itoa(deletes))
	_, _ = h.out.WriteString("</h3>\n")
	_, _ = h.out.WriteString(`<p>`)
	//_, _ = h.out.WriteString(diffMatch.DiffPrettyHtml(diffs))
	_, _ = h.out.WriteString(diffHtml)
	_, _ = h.out.WriteString("</p>\n")
}

func (h *HTMLWriter) WriteEnd(insertSum int, deleteSum int, diffCount int) {
	//fmt.Println("Num Diff", diffCount)
	_, _ = h.out.WriteString(`<p>Total Inserted Chars `)
	_, _ = h.out.WriteString(strconv.Itoa(insertSum))
	_, _ = h.out.WriteString(`, Total Deleted Chars `)
	_, _ = h.out.WriteString(strconv.Itoa(deleteSum))
	_, _ = h.out.WriteString("</p>\n")
	_, _ = h.out.WriteString(`<p>`)
	_, _ = h.out.WriteString("Total Difference Count: ")
	_, _ = h.out.WriteString(strconv.Itoa(diffCount))
	_, _ = h.out.WriteString("</p>\n")
	end := ` </body>
</html>`
	_, _ = h.out.WriteString(end)
	_ = h.out.Close()
}
