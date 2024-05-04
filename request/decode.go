package request

import (
	"bytes"
	"context"
	"dataset"
	log "dataset/logger"
	"gopkg.in/yaml.v3"
)

type RequestDecoder struct {
	ctx context.Context
}

func NewRequestDecoder(ctx context.Context) RequestDecoder {
	var r RequestDecoder
	r.ctx = ctx
	return r
}

func (r *RequestDecoder) Process(yamlRequest []byte) (Request, dataset.Status) {
	var request Request
	var status dataset.Status
	request, status = r.Decode(yamlRequest)
	if status.IsErr {
		return request, status
	}
	status = r.Validate(&request)
	if status.IsErr {
		return request, status
	}
	r.Prereq(&request)
	status = r.Depend(request)
	return request, status
}

func (r *RequestDecoder) Decode(requestYaml []byte) (Request, dataset.Status) {
	var resp Request
	var status dataset.Status
	reader := bytes.NewReader(requestYaml)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	err := decoder.Decode(&resp)
	if err != nil {
		status = log.Error(r.ctx, 500, err, `Error decoding YAML to request`)
	}
	resp.Testament.BuildBookMaps() // Builds Map for t.HasOT(bookId), t.HasNT(bookId)
	return resp, status
}

func (r *RequestDecoder) Encode(req Request) (string, dataset.Status) {
	var status dataset.Status
	d, err := yaml.Marshal(&req)
	if err != nil {
		status = log.Error(r.ctx, 500, err, `Error encoding request to YAML`)
		return ``, status
	}
	return string(d), status
}
