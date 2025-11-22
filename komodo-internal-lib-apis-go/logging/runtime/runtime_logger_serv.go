package runtime

import (
	"encoding/json"
	"fmt"
	"komodo-internal-lib-apis-go/config"
	"log"
	"net/http"
	"strings"
	"sync"
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

type LoggerConfig struct {
	EnableRemoteLogs bool
	Host             string
	LogLevel         string
}

var (
	enableRemoteLogs bool
	logQueue         chan logPayload
	host             string
	logLevel         string
	initOnce         sync.Once
)

// ensureInitialized lazily initializes the logger on first use
func ensureInitialized() {
	initOnce.Do(func() {
		// Read config values (now they should be available)
		enableRemoteLogs = config.GetConfigValue("ENABLE_REMOTE_LOGS") == "true"
		host = config.GetConfigValue("HOSTNAME")
		logLevel = strings.ToUpper(config.GetConfigValue("LOG_LEVEL"))
		
		if logLevel == "" {
			logLevel = "INFO" // Default to INFO, not ERROR
		}

		logQueue = make(chan logPayload, 1000)

		// Start the remote log consumer goroutine
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

		log.Println("[INFO] Logger initialized:", "level="+logLevel, "remote=", enableRemoteLogs)
	})
}

func logByLevel(level, msg string, meta any) {
	ensureInitialized()

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

func sanitizeMeta(meta ...any) any {
	if len(meta) == 0 { return nil }

	switch m := meta[0].(type) {
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

func InitLogger(config ...LoggerConfig) {
	ensureInitialized() 

	if len(config) > 0 {
		cfg := config[0]
		enableRemoteLogs = cfg.EnableRemoteLogs
		host = cfg.Host
		if cfg.LogLevel != "" {
			logLevel = strings.ToUpper(cfg.LogLevel)
		}
	}
}

// ===== All-purpose logging functions =====

func Info(msg string, details ...any)  {
	m := sanitizeMeta(details...)
	logByLevel("INFO", msg, m)
	enqueueRemoteLog("INFO", msg, m)
}
func Warn(msg string, details ...any)  {
	m := sanitizeMeta(details...)
	logByLevel("WARN", msg, m)
	enqueueRemoteLog("WARN", msg, m)
}
func Error(msg string, details ...any) {
	m := sanitizeMeta(details...)
	logByLevel("ERROR", msg, m)
	enqueueRemoteLog("ERROR", msg, m)
}
func Fatal(msg string, details ...any) {
	m := sanitizeMeta(details...)
	logByLevel("FATAL", msg, m)
	enqueueRemoteLog("FATAL", msg, m)
}
func Debug(msg string, details ...any) {
	m := sanitizeMeta(details...)
	logByLevel("DEBUG", msg, m)
	enqueueRemoteLog("DEBUG", msg, m)
}
func Trace(msg string, details ...any) {
	m := sanitizeMeta(details...)
	logByLevel("TRACE", msg, m)
	enqueueRemoteLog("TRACE", msg, m)
}

// ===== Prints logs only to stdout ===== 

func Print(level string, msg string, details ...any) { logByLevel(level, msg, sanitizeMeta(details)) }
func PrintFatal(msg string, details ...any) { logByLevel("FATAL", msg, sanitizeMeta(details)) }
func PrintError(msg string, details ...any) { logByLevel("ERROR", msg, sanitizeMeta(details)) }
func PrintWarn(msg string, details ...any)  { logByLevel("WARN", msg, sanitizeMeta(details)) }
func PrintInfo(msg string, details ...any)  { logByLevel("INFO", msg, sanitizeMeta(details)) }
func PrintDebug(msg string, details ...any) { logByLevel("DEBUG", msg, sanitizeMeta(details)) }
func PrintTrace(msg string, details ...any) { logByLevel("TRACE", msg, sanitizeMeta(details)) }
