package intercept

import (
	"net/http"

	"allseer/internal/rules"
)

type ResponseInterceptor interface {
	AfterResponse(resp *http.Response, decision rules.Decision) error
}
