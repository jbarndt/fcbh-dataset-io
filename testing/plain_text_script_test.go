package testing

import (
	"fmt"
	"testing"
)

const PlainTextScript = `Required:
  IsNew: yes
  RequestName: PlainTextScript
  BibleId: ENGWEB
TextData:
  BibleBrain:
    TextPlain: yes
OutputFormat:
  CSV: yes
`

func TestPlainTextScript(t *testing.T) {
	csvResp, status := HttpPost(`PlainTextScript.csv`, PlainTextScript, t)
	fmt.Printf("Response status: %d\n", status)
	fmt.Printf("Response body: %s\n", string(csvResp))
	if len(csvResp) < 1000 {
		t.Error(`Expected at least 1000 rows, got`, len(csvResp))
	}
}
