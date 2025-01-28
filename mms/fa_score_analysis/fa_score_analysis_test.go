package fa_score_analysis

import (
	"context"
	"dataset/db"
	"dataset/fetch"
	"fmt"
	"testing"
)

func TestFAScoreAnalysis(t *testing.T) {
	ctx := context.Background()
	user, _ := fetch.GetTestUser()
	list := []string{"N2KTB_ESB", "N2CFM_BSM", "N2CHF_TBL", "N2CUL_MNT"}
	for _, database := range list {
		conn, status := db.NewerDBAdapter(ctx, false, user.Username, database)
		if status != nil {
			t.Fatal(status)
		}
		output, status := FAScoreAnalysis(conn)
		if status != nil {
			t.Fatal(status)
		}
		conn.Close()
		fmt.Println(output)
	}
}
