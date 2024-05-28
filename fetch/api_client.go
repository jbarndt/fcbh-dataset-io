package fetch

import (
	"context"
	"dataset"
	log "dataset/logger"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	HOST = "https://4.dbt.io/api/"
)

func HttpGet(ctx context.Context, url string, desc string) ([]byte, dataset.Status) {
	return httpGet(ctx, url, false, desc)
}

func httpGet(ctx context.Context, url string, ok403 bool, desc string) ([]byte, dataset.Status) {
	var body []byte
	var status dataset.Status
	if strings.Contains(url, HOST) {
		url += `&limit=100000&key=` + os.Getenv(`FCBH_DBP_KEY`)
	}
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		status = log.Error(ctx, resp.StatusCode, err, "Error in DBP API request for:", desc)
		return body, status
	}
	if ok403 && resp.StatusCode == 403 {
		status.IsErr = true
		status.Status = 403
		return body, status
	}
	if resp.Status[0] != '2' {
		status = log.ErrorNoErr(ctx, resp.StatusCode, resp.Status, desc)
		return body, status
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		status = log.Error(ctx, resp.StatusCode, err, "Error reading DBP API response for:", desc)
	}
	return body, status
}
