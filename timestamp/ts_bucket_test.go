package timestamp

import (
	"context"
	"fmt"
	"testing"
)

// This test is not working...
func TestTSBucket(t *testing.T) {
	ctx := context.Background()
	ts, status := NewTSBucket(ctx)
	if status.IsErr {
		t.Fatal(status)
	}
	key, status := ts.GetKey(ScriptTS, `ENGWEBN2DA`, `REV`, 22)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(key)
	object, status := ts.GetObject(TSBucketName, key)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(string(object))
	timestamps, status := ts.GetTimestamps(ScriptTS, `ENGWEBN2DA`, `REV`, 22)
	if status.IsErr {
		t.Fatal(status)
	}
	for _, time := range timestamps {
		fmt.Println(time)
	}
}

func TestTSBucket_GetTimestamps(t *testing.T) {
	ctx := context.Background()
	ts, status := NewTSBucket(ctx)
	if status.IsErr {
		t.Fatal(status)
	}
	key, status := ts.GetKey(VerseAeneas, `ENGWEBN2DA`, `REV`, 22)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(key)
	object, status := ts.GetObject(TSBucketName, key)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(string(object))
	timestamps, status := ts.GetTimestamps(VerseAeneas, `ENGWEBN2DA`, `REV`, 22)
	if status.IsErr {
		t.Fatal(status)
	}
	for _, time := range timestamps {
		fmt.Println(time)
	}
}

//aws s3 ls s3://dbp-aeneas-staging/Latin_N2_organized/pass_qc/ENGWEBN2DA/cue_info_text

//aeneas_verse_timings

//aws s3 ls s3://dbp-aeneas-staging/Latin_N2_organized/pass_qc/ENGWEBN2DA/aeneas_verse_timings/
