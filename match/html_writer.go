package match

import (
	"context"
	"dataset"
	log "dataset/logger"
	"math"
	"os"
	"strconv"
	"time"
)

type HTMLWriter struct {
	ctx         context.Context
	datasetName string
	out         *os.File
	lineNum     int
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

func (h *HTMLWriter) WriteHeading(baseDataset string) string {
	head := `<!DOCTYPE html>
<html>
 <head>
  <meta charset="utf-8">
  <title>File Difference</title>
`
	_, _ = h.out.WriteString(head)
	_, _ = h.out.WriteString(`<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/1.10.21/css/jquery.dataTables.css">`)
	_, _ = h.out.WriteString("</head><body>\n")
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
	checkbox := `<div style="text-align: center; margin: 10px;">
		<input type="checkbox" id="hideVerse0"><label for="hideVerse0">Hide Headings</label>
	</div>
`
	_, _ = h.out.WriteString(checkbox)
	table := `<table id="diffTable" class="display">
    <thead>
    <tr>
        <th>Line</th>
        <th>Err %</th>
		<th>Chars</th>
		<th>Inbal</th>
		<th>Error</th>
        <th>Ref</th>
		<th>Text Comparison</th>
    </tr>
    </thead>
    <tbody>
`
	_, _ = h.out.WriteString(table)
	return h.out.Name()
}

func (h *HTMLWriter) WriteVerseDiff(verse pair, inserts int, deletes int, errPct float64, diffHtml string) {
	h.lineNum++
	_, _ = h.out.WriteString("<tr>\n")
	h.writeCell(strconv.Itoa(h.lineNum))
	h.writeCell(strconv.FormatFloat(errPct, 'f', 0, 64))
	h.writeCell(strconv.Itoa(inserts + deletes))
	h.writeCell(strconv.Itoa(int(math.Abs(float64(inserts - deletes)))))
	h.writeCell(`+` + strconv.Itoa(inserts) + ` -` + strconv.Itoa(deletes))
	h.writeCell(verse.bookId + ` ` + strconv.Itoa(verse.chapter) + `:` + verse.num)
	h.writeCell(diffHtml)
	_, _ = h.out.WriteString("</tr>\n")
}

func (h *HTMLWriter) WriteChapterDiff(bookId string, chapter int, inserts int, deletes int, errPct float64, diffHtml string) {
	h.lineNum++
	_, _ = h.out.WriteString("<tr>\n")
	h.writeCell(strconv.Itoa(h.lineNum))
	h.writeCell(strconv.FormatFloat(errPct, 'f', 0, 64))
	h.writeCell(strconv.Itoa(inserts + deletes))
	h.writeCell(strconv.Itoa(int(math.Abs(float64(inserts - deletes)))))
	h.writeCell(`+` + strconv.Itoa(inserts) + ` -` + strconv.Itoa(deletes))
	h.writeCell(bookId + ` ` + strconv.Itoa(chapter))
	h.writeCell(diffHtml)
	_, _ = h.out.WriteString("</tr>\n")
}

func (h *HTMLWriter) WriteScriptLineDiff(bookId string, chapter int, line string, inserts int, deletes int, errPct float64, diffHtml string) {
	_, _ = h.out.WriteString("<tr>\n")
	h.writeCell(line)
	h.writeCell(strconv.FormatFloat(errPct, 'f', 0, 64))
	h.writeCell(strconv.Itoa(inserts + deletes))
	h.writeCell(strconv.Itoa(int(math.Abs(float64(inserts - deletes)))))
	h.writeCell(`+` + strconv.Itoa(inserts) + ` -` + strconv.Itoa(deletes))
	h.writeCell(bookId + ` ` + strconv.Itoa(chapter))
	h.writeCell(diffHtml)
	_, _ = h.out.WriteString("</tr>\n")
}

func (h *HTMLWriter) writeCell(content string) {
	_, _ = h.out.WriteString(`<td>`)
	_, _ = h.out.WriteString(content)
	_, _ = h.out.WriteString(`</td>`)
}

func (h *HTMLWriter) WriteEnd(insertSum int, deleteSum int, diffCount int) {
	table := `</tbody>
	</table>
`
	_, _ = h.out.WriteString(table)
	_, _ = h.out.WriteString(`<p>Total Inserted Chars `)
	_, _ = h.out.WriteString(strconv.Itoa(insertSum))
	_, _ = h.out.WriteString(`, Total Deleted Chars `)
	_, _ = h.out.WriteString(strconv.Itoa(deleteSum))
	_, _ = h.out.WriteString("</p>\n")
	_, _ = h.out.WriteString(`<p>`)
	_, _ = h.out.WriteString("Total Difference Count: ")
	_, _ = h.out.WriteString(strconv.Itoa(diffCount))
	_, _ = h.out.WriteString("</p>\n")
	_, _ = h.out.WriteString(`<script type="text/javascript" src="https://code.jquery.com/jquery-3.5.1.js"></script>`)
	_, _ = h.out.WriteString("\n")
	_, _ = h.out.WriteString(`<script type="text/javascript" src="https://cdn.datatables.net/1.10.21/js/jquery.dataTables.js"></script>`)
	_, _ = h.out.WriteString("\n")
	style := `<style>
	.dataTables_length select {
	width: auto;
	display: inline-block;
	padding: 5px;
		margin-left: 5px;
		border-radius: 4px;
	border: 1px solid #ccc;
	}
	.dataTables_filter input {
		width: auto;
		display: inline-block;
		padding: 5px;
		border-radius: 4px;
		border: 1px solid #ccc;
	}
	.dataTables_wrapper .dataTables_length, .dataTables_wrapper .dataTables_filter {
		margin-bottom: 20px;
	}
	</style>
`
	_, _ = h.out.WriteString(style)
	script := `<script>
    $(document).ready(function() {
        var table = $('#diffTable').DataTable({
            "columnDefs": [
                { "orderable": false, "targets": [4,6] }
				// { "visible": false, "targets": [7] }  
            ],
            "pageLength": 10,
            "lengthMenu": [[5, 10, 25, 50, -1], [5, 10, 25, 50, "All"]]
        });
        $('#hideVerse0').on('change', function() {
            if(this.checked) {
                table.rows(function(idx, data, node) {
					return data[5].endsWith(":0");
                }).nodes().to$().hide();
            } else {
                table.rows().nodes().to$().show();
            }
            // Optional: redraw table to fix any layout issues
            // table.draw();
        });
    });
</script>
</body>
</html>
`
	_, _ = h.out.WriteString(script)
	_ = h.out.Close()
}
