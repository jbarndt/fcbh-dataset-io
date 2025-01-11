package fa_score_analysis

import (
	"context"
	"dataset/db"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFAScoreAnalysis(t *testing.T) {
	ctx := context.Background()
	for _, filename := range []string{"N2KTB_ESB.db", "N2CFM_BSM.db", "N2CHF_TBL.db"} {
		database := filepath.Join(os.Getenv("HOME"), filename)
		conn := db.NewDBAdapter(ctx, database)
		output, status := FAScoreAnalysis(conn)
		if status.IsErr {
			t.Fatal(status)
		}
		conn.Close()
		fmt.Println(output)
	}
}
