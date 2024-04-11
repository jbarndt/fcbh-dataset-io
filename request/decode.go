package request

import (
	"context"
	"dataset"
	log "dataset/logger"
	"gopkg.in/yaml.v3"
	"os"
)

type RequestDecoder struct {
	ctx context.Context
}

func NewRequestDecoder(ctx context.Context) RequestDecoder {
	var r RequestDecoder
	r.ctx = ctx
	return r
}

func (r *RequestDecoder) Process(str string) dataset.Status {
	var request Request
	var status dataset.Status
	if len(str) < 50 {
		request, status = r.DecodeFile(str)
	} else {
		request, status = r.DecodeString(str)
	}
	if status.IsErr {
		return status
	}
	status = r.Validate(request)
	if status.IsErr {
		return status
	}
	r.Prereq(&request)
	status = r.Depend(request)
	if status.IsErr {
		return status
	}
	return status
}

func (r *RequestDecoder) DecodeFile(path string) (Request, dataset.Status) {
	var resp Request
	var status dataset.Status
	content, err := os.ReadFile(path)
	if err != nil {
		status = log.Error(r.ctx, 500, err, `Error reading YAML file`)
		return resp, status
	}
	return r.decode(content)
}

func (r *RequestDecoder) DecodeString(str string) (Request, dataset.Status) {
	return r.decode([]byte(str))
}

func (r *RequestDecoder) decode(requestYaml []byte) (Request, dataset.Status) {
	var resp Request
	var status dataset.Status
	err := yaml.Unmarshal(requestYaml, &resp)
	if err != nil {
		status = log.Error(r.ctx, 500, err, `Error decoding YAML to request`)
	}
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
