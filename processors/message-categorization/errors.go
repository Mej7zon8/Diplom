package message_categorization

import "github.com/go-resty/resty/v2"

type UnexpectedResponse struct {
	StatusCode int
	Status     string
	Body       string
}

func newUnexpectedResponse(r *resty.Response) *UnexpectedResponse {
	return &UnexpectedResponse{
		StatusCode: r.StatusCode(),
		Status:     r.Status(),
		Body:       r.String(),
	}
}

func (u *UnexpectedResponse) Error() string {
	var res = "unexpected response: " + u.Status
	if u.Body != "" && u.Body[0] == '{' {
		res += ": " + u.Body
	}
	return res
}
