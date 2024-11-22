package run_control

import (
	"context"
	"dataset/db"
	"os"
	"testing"
)

const runBucketTest = `is_new: yes
dataset_name: MyProject
bible_id: ENGWEB
username: GaryNTest
email: gary@shortsands.com
output_file: abc/my_project.csv
`

func TestRunBucket(t *testing.T) {
	ctx := context.Background()
	b := NewRunBucket(ctx)
	b.isUnitTest = true
	b.AddYaml([]byte(runBucketTest))
	if b.username != "GaryNTest" {
		t.Error("Username should be GaryNTest, it is: ", b.username)
	}
	if len(b.username) != 9 {
		t.Error("Username should be 9 characters")
	}
	if b.dataset != "MyProject" {
		t.Error("Project should be MyProject, it is:", b.dataset)
	}
	b.AddLogFile(os.Getenv("FCBH_DATASET_LOG_FILE"))
	database1, status := db.NewerDBAdapter(ctx, true, b.username, "TestRunBucket1")
	if status.IsErr {
		t.Fatal(status.String())
	}
	b.AddDatabase(database1)
	database2, status := db.NewerDBAdapter(ctx, true, b.username, "TestRunBucket2")
	if status.IsErr {
		t.Fatal(status.String())
	}
	b.AddDatabase(database2)
	b.AddOutput("../tests/02__plain_text_edit_script.csv")
	b.AddOutput("../tests/02__plain_text_edit_script.json")
	b.AddDuration()
	status = b.PersistToBucket()
	if status.IsErr {
		t.Fatal(status.String())
	}
}
