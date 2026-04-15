package intercept

import (
	"net/http"

	"allseer/internal/rules"
)

type RequestInterceptor interface {
	BeforeForward(req *http.Request, decision rules.Decision) error
}
