package message_categorization

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"messenger/data/entities"
	"messenger/data/providers/advertisement"
	"strings"
)

const baseUrl = "http://localhost:9000"

func resolveUrl(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return baseUrl + path
}

var Instance = New()

// The Processor is used for categorizing messages.
// The algorithm is implemented in the python service 'messenger-categorization-service'.
type Processor struct {
	client *resty.Client
}

func New() *Processor {
	return &Processor{
		client: resty.New(),
	}
}

type processReq struct {
	Text            string                     `json:"text"`
	CandidateLabels []entities.AdCategoryLabel `json:"candidate_labels"`
}

type ProcessRes struct {
	Result map[string]float64 `json:"result"`
}

func (p *Processor) Process(ctx context.Context, message string) (result ProcessRes, e error) {
	return genericDo[ProcessRes](p, ctx, "/categorize", processReq{
		Text:            message,
		CandidateLabels: advertisement.Instance.ListLabels(),
	}, false)
}

func genericDo[T any](p *Processor, ctx context.Context, path string, body any, skipUnmarshal bool) (T, error) {
	var req = p.client.R().SetContext(ctx)
	if body != nil {
		req.SetBody(body)
	}

	r, e := req.Post(resolveUrl(path))
	var value T
	if e != nil {
		return value, e
	}
	if r.IsError() {
		return value, newUnexpectedResponse(r)
	}
	if skipUnmarshal {
		return value, nil
	}
	data := r.Body()

	e = json.Unmarshal(data, &value)
	return value, e
}
