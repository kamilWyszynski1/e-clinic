package http_client

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

func ToRetryableReq(req *http.Request) (*retryablehttp.Request, error) {
	rreq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, fmt.Errorf("cant convert to retryable req %w", err)
	}

	return rreq, nil
}
