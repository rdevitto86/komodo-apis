package latency

import (
	ctxKeys "komodo-internal-lib-apis-go/common/context"
	"net/http"
	"time"
)

func SetLatencyHeaders(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set(string(ctxKeys.EndTimeKey), time.Now().UTC().Format(time.RFC3339))

	if start, ok := req.Context().Value(ctxKeys.StartTimeKey).(time.Time); ok {
		duration := time.Since(start)
		wtr.Header().Set(string(ctxKeys.DurationKey), duration.String())
	}
}