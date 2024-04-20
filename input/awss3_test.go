package input

import (
	"context"
	"dataset/request"
	"os"
	"path/filepath"
	"testing"
)

func TestAWSS3Audio(t *testing.T) {
	ctx := context.Background()
	key := `s3://dbp-prod/audio/ENGWEB/ENGWEBN2DA/*.mp3`
	testament := request.Testament{NT: true}
	testament.BuildBookMaps()
	//directory := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), `ENGWEB`, `ENGWEBN2DA`)
	//os.Remove(filepath.Join(directory, `040MAT.usx`))
	//os.Remove(filepath.Join(directory, `041MRK.usx`))
	//os.Remove(filepath.Join(directory, `042LUK.usx`))
	//os.Remove(filepath.Join(directory, `043JHN.usx`))
	files, status := AWSS3Input(ctx, key, testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	if len(files) != 260 {
		t.Error(`len(files) should be 260`, len(files))
	}
}

func TestAWSS3USX(t *testing.T) {
	ctx := context.Background()
	key := `s3://dbp-prod/text/ENGWEB/ENGWEBN_ET-usx`
	testament := request.Testament{NT: true}
	testament.BuildBookMaps()
	directory := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), `ENGWEB`, `ENGWEBN_ET-usx`)
	os.Remove(filepath.Join(directory, `040MAT.usx`))
	os.Remove(filepath.Join(directory, `041MRK.usx`))
	os.Remove(filepath.Join(directory, `042LUK.usx`))
	os.Remove(filepath.Join(directory, `043JHN.usx`))
	files, status := AWSS3Input(ctx, key, testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	if len(files) != 27 {
		t.Error(`len(files) should be 27`, len(files))
	}
	if files[0].MediaId != `ENGWEBN_ET-usx` {
		t.Error(`files[0].MediaId should be ENGWEBN_ET-usx`, files[0].MediaId)
	}
	if files[0].MediaType != `text_usx` {
		t.Error(`files[0].MediaType should be text_usx`, files[0].MediaType)
	}
}

func TestAWSS3Zip(t *testing.T) {
	ctx := context.Background()
	key := `s3://dbp-prod/audio/ENGWEB/ENGWEBN2DA/ENGWEBN2DA.zip`
	testament := request.Testament{NT: true}
	testament.BuildBookMaps()
	//directory := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), `ENGWEB`, `ENGWEBN2DA`)
	//os.Remove(filepath.Join(directory, `040MAT.usx`))
	//os.Remove(filepath.Join(directory, `041MRK.usx`))
	//os.Remove(filepath.Join(directory, `042LUK.usx`))
	//os.Remove(filepath.Join(directory, `043JHN.usx`))
	files, status := AWSS3Input(ctx, key, testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	if len(files) != 1 {
		t.Error(`len(files) should be 260`, len(files))
	}
}
