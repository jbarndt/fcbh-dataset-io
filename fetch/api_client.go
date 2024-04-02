package fetch

import (
	"context"
	"dataset"
	log "dataset/logger"
	"io"
	"net/http"
	"os"
)

const (
	HOST = "https://4.dbt.io/api/"
)

func httpGet(ctx context.Context, url string, desc string) ([]byte, dataset.Status) {
	var body []byte
	var status dataset.Status
	url += `&limit=100000&key=` + os.Getenv(`FCBH_DBP_KEY`)
	resp, err := http.Get(url)
	if err != nil {
		status = log.Error(ctx, resp.StatusCode, err, "Error in DBP API request for:", desc)
		return body, status
	}
	defer resp.Body.Close()
	if resp.Status[0] == '2' {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			status = log.Error(ctx, resp.StatusCode, err, "Error reading DBP API response for:", desc)
			return body, status
		}
	}
	return body, status
}
