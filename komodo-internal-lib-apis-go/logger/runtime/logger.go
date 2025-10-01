package runtime

import (
	"encoding/json"
	"fmt"
	"komodo-internal-lib-apis-go/config"
	"log"
	"net/http"
	"strings"
	"time"
)

var logPriorities = map[string]int{
	"DEBUG": 0,
	"debug": 0,
	"TRACE": 0,
	"trace": 0,
	"INFO":  1,
	"info":  1,
	"WARN":  2,
	"warn":  2,
	"ERROR": 3,
	"error": 3,
	"FATAL": 4,
	"fatal": 4,
}

// payload sent to external backends
type logPayload struct {
	Level     string      `json:"level"`
	Msg       string      `json:"msg"`
	Meta      any         `json:"meta,omitempty"`
	Timestamp time.Time   `json:"ts"`
	Host      string      `json:"host,omitempty"`
}

var (
	enableRemoteLogs    = config.GetConfigValue("ENABLE_REMOTE_LOGS") == "true"
	logQueue         		= make(chan logPayload, 1000) // bounded queue of JSON payloads
	host             		= config.GetConfigValue("HOSTNAME")
	logLevel        		= config.GetConfigValue("LOG_LEVEL")
)

func init() {
	if lvl := strings.ToUpper(logLevel); lvl != "" {
		logLevel = lvl
	}
	if logLevel == "" {
		logLevel = "ERROR"
	}

	if enableRemoteLogs {
		go func() {
			for pl := range logQueue {
				b, err := json.Marshal(pl)
				if err != nil {
					fmt.Printf("[ERROR] failed to marshal log payload: %v\n", err)
					continue
				}
				if err := postRemoteLog(b); err != nil {
					fmt.Printf("[ERROR] failed to send log to remote: %v\n", err)
				}
			}
		}()
	}
}

func logByLevel(level, msg string, meta any) {
	// ignore non-scoped logs
	if logPriorities[level] < logPriorities[logLevel] {
		return
	}

	if meta != nil {
		log.Println("[" + level + "]", msg, meta)
	} else {
		log.Println("[" + level + "]", msg)
	}
}

func enqueueRemoteLog(level string, msg string, meta any) error {
	if !enableRemoteLogs || logPriorities[level] < logPriorities[logLevel] {
		return nil
	}

	payload := logPayload{
		Level:     level,
		Msg:       msg,
		Meta:      meta,
		Timestamp: time.Now().UTC(),
		Host:      host,
	}

	select {
		case logQueue <- payload:
		case <-time.After(50 * time.Millisecond):
			// queue full — drop silently
	}
	return nil
}

func postRemoteLog(raw []byte) error {
	// TODO: push to Grafana Loki or other Grafana datasource.
	_ = raw
	return nil
}

func sanitizeMeta(meta any) any {
	switch m := meta.(type) {			
		case string, map[string]any:
			return m
		case error:
			return m.Error()
		case *time.Time:
			return m.UTC().Format(time.RFC3339Nano)
		case *http.Request:
			return map[string]interface{}{
				"method":     	m.Method,
				"path":     		m.URL.Path,
				"host":     		m.Host,
				"request_id": 	m.Context().Value("request_id"),
				"pathParams": 	m.Context().Value("pathParams"),
				"queryParams": 	m.URL.Query(),
				"headers":   		m.Header,
				"remoteAddr": 	m.RemoteAddr,
			}
		case json.RawMessage:
			var tmp any
			if err := json.Unmarshal(m, &tmp); err == nil {
				return tmp
			}
			return nil
		case []byte:
			var tmp any
			if err := json.Unmarshal(m, &tmp); err == nil {
				return tmp
			}
			return nil
	}
	// TODO - support arrays
	return nil // unsupported type
}

func Info(msg string, meta any)  {
	meta = sanitizeMeta(meta)
	logByLevel("INFO", msg, meta)
	enqueueRemoteLog("INFO", msg, meta)
}
func Warn(msg string, meta any)  {
	meta = sanitizeMeta(meta)
	logByLevel("WARN", msg, meta)
	enqueueRemoteLog("WARN", msg, meta)
}
func Error(msg string, meta any) {
	meta = sanitizeMeta(meta)
	logByLevel("ERROR", msg, meta)
	enqueueRemoteLog("ERROR", msg, meta)
}
func Fatal(msg string, meta any) {
	meta = sanitizeMeta(meta)
	logByLevel("FATAL", msg, meta)
	enqueueRemoteLog("FATAL", msg, meta)
}
func Debug(msg string, meta any) {
	meta = sanitizeMeta(meta)
	logByLevel("DEBUG", msg, meta)
	enqueueRemoteLog("DEBUG", msg, meta)
}
func Trace(msg string, meta any) {
	meta = sanitizeMeta(meta)
	logByLevel("TRACE", msg, meta)
	enqueueRemoteLog("TRACE", msg, meta)
}

// Prints logs only to stdout
func Print(level string, msg string, meta any) { logByLevel(level, msg, sanitizeMeta(meta)) }
