package latency

import (
	ctxKeys "komodo-internal-lib-apis-go/common/context"
	"net/http"
	"time"
)

func SetLatencyHeaders(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set(string(ctxKeys.END_TIME_KEY), time.Now().UTC().Format(time.RFC3339))

	if start, ok := req.Context().Value(ctxKeys.START_TIME_KEY).(time.Time); ok {
		duration := time.Since(start)
		wtr.Header().Set(string(ctxKeys.DURATION_KEY), duration.String())
	}
}