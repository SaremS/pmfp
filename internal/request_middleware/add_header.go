package request_middleware

import (
	"net/http"
)

type AddHeader struct {
	HeaderName  string `json:"header_name"`
	HeaderValue string `json:"header_value"`
}

func (a *AddHeader) Apply(request *http.Request) {
	request.Header.Set(a.HeaderName, a.HeaderValue)
}

func init() {
	RegisterRequestMiddleware("add_header", func() RequestMiddleware { return &AddHeader{} })
}
