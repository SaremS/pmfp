package request_manipulation

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
	RegisterRequestManipulator("add_header", func() RequestManipulator { return &AddHeader{} })
}
