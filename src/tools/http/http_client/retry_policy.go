package http_client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

//GetStatusCodeRetryPolicy returns retry policy that calls retry if status is in one of the provided ranges
func GetStatusCodeRetryPolicy(ranges []*StatusCodeRange, log logrus.FieldLogger) retryablehttp.CheckRetry {
	return func(_ context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			return true, err
		}

		for _, rng := range ranges {
			if resp.StatusCode >= rng.Min && resp.StatusCode <= rng.Max {
				if log != nil {
					log.Warning(fmt.Sprintf("unexpected code %d on %s", resp.StatusCode, resp.Request.URL))
				}
				return true, nil
			}
		}

		return false, nil
	}
}
