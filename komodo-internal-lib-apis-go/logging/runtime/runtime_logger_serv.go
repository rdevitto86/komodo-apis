package runtime

import (
	"encoding/json"
	"komodo-internal-lib-apis-go/config"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ===== Configuration =====

type LoggerConfig struct {
	AppName          string
	EnableRemoteLogs bool
	Host             string
	LogLevel         string
}

func InitLogger(cfg ...LoggerConfig) {
	initOnce.Do(func() { initialize(cfg...) })
}

// ===== Public Logging Functions =====

func Info(msg string, details ...any)  { writeLog("INFO", msg, true, details...) }
func Warn(msg string, details ...any)  { writeLog("WARN", msg, true, details...) }
func Error(msg string, details ...any) { writeLog("ERROR", msg, true, details...) }
func Fatal(msg string, details ...any) { writeLog("FATAL", msg, true, details...) }
func Debug(msg string, details ...any) { writeLog("DEBUG", msg, true, details...) }
func Trace(msg string, details ...any) { writeLog("TRACE", msg, true, details...) }

// ===== Stdout-Only Functions (no remote logging) =====

func Print(level string, msg string, details ...any) { writeLog(level, msg, false, details...) }
func PrintFatal(msg string, details ...any) { writeLog("FATAL", msg, false, details...) }
func PrintError(msg string, details ...any) { writeLog("ERROR", msg, false, details...) }
func PrintWarn(msg string, details ...any)  { writeLog("WARN", msg, false, details...) }
func PrintInfo(msg string, details ...any)  { writeLog("INFO", msg, false, details...) }
func PrintDebug(msg string, details ...any) { writeLog("DEBUG", msg, false, details...) }
func PrintTrace(msg string, details ...any) { writeLog("TRACE", msg, false, details...) }

// ===== Internal State =====

var (
	appName						string
	enableRemoteLogs 	bool
	logQueue         	chan logPayload
	host             	string
	logLevel         	string
	initOnce         	sync.Once
)

var logPriorities = map[string]int{
	"TRACE": 0,
	"DEBUG": 1,
	"INFO":  2,
	"WARN":  3,
	"ERROR": 4,
	"FATAL": 5,
}

type logPayload struct {
	AppName   string    `json:"app,omitempty"`
	Level     string    `json:"level"`
	Msg       string    `json:"msg"`
	Meta      any       `json:"meta,omitempty"`
	Timestamp time.Time `json:"ts"`
	Host      string    `json:"host,omitempty"`
}

// ===== Initialization =====

func initialize(cfg ...LoggerConfig) {
	log.Println("[INFO] initializing logger...")

	// Load from config (if provided)
	if len(cfg) > 0 {
		c := cfg[0]
		appName = c.AppName
		enableRemoteLogs = c.EnableRemoteLogs
		host = c.Host
		logLevel = strings.ToUpper(c.LogLevel)
	}

	// Fallback to environment variables if not set via config
	if appName == "" { appName = config.GetConfigValue("APP_NAME") }
	if logLevel == "" {
		if strings.ToLower(config.GetConfigValue("ENV")) == "local" {
			logLevel = "DEBUG"
		} else if logLevel = strings.ToUpper(config.GetConfigValue("LOG_LEVEL")); logLevel == "" {
			logLevel = "ERROR"
		}
	}
	if host == "" { host = config.GetConfigValue("LOG_HOSTNAME") }
	if len(cfg) == 0 { enableRemoteLogs = config.GetConfigValue("LOG_ENABLE_REMOTE") == "true" }
	
	logQueue = make(chan logPayload, 1000)

	// Start remote log consumer if enabled
	if enableRemoteLogs {
		go consumeRemoteLogs()
	}

	log.Printf(
		"[INFO] logger initialized (app=%s, level=%s, remote=%v, host=%s)\n",
		appName, logLevel, enableRemoteLogs, host,
	)
}

func consumeRemoteLogs() {
	for payload := range logQueue {
		data, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[ERROR] failed to marshal log payload: %v", err)
			continue
		}
		if err := sendRemoteLog(data); err != nil {
			log.Printf("[ERROR] failed to send log to remote: %v", err)
		}
	}
}

// ===== Core Logging Logic =====

func writeLog(level, msg string, remote bool, details ...any) {
	initOnce.Do(func() { initialize() })

	// Filter by log level
	if logPriorities[level] < logPriorities[logLevel] {
		return
	}

	meta := sanitizeMeta(details...)

	// Print to stdout
	if meta != nil {
		log.Printf("[%s] %s %v\n", level, msg, meta)
	} else {
		log.Printf("[%s] %s\n", level, msg)
	}

	// Enqueue for remote logging
	if remote && enableRemoteLogs {
		select {
			case logQueue <- logPayload{
				AppName:   appName,
				Level:     level,
				Msg:       msg,
				Meta:      meta,
				Timestamp: time.Now().UTC(),
				Host:      host,
			}:
			case <-time.After(50 * time.Millisecond):
				// Queue full, drop silently
		}
	}
}

func sendRemoteLog(data []byte) error {
	// TODO: push to Splunk, New Relic, Grafana Loki, or custom endpoint
	_ = data
	return nil
}

// ===== Metadata Sanitization =====

func sanitizeMeta(details ...any) any {
	if len(details) == 0 { return nil }

	switch meta := details[0].(type) {
		case nil:
			return nil
		case string, map[string]any:
			return meta
		case error:
			return meta.Error()
		case *time.Time:
			if meta != nil {
				return meta.UTC().Format(time.RFC3339Nano)
			}
			return nil
		case time.Time:
			return meta.UTC().Format(time.RFC3339Nano)
		case *http.Request:
			return map[string]any{
				"method":      meta.Method,
				"path":        meta.URL.Path,
				"host":        meta.Host,
				"request_id":  meta.Context().Value("request_id"),
				"pathParams":  meta.Context().Value("pathParams"),
				"queryParams": meta.URL.Query(),
				"headers":     meta.Header,
				"remoteAddr":  meta.RemoteAddr,
			}
		case json.RawMessage, []byte:
			var tmp any
			var data []byte
			if raw, ok := meta.(json.RawMessage); ok {
				data = raw
			} else {
				data = meta.([]byte)
			}
			if err := json.Unmarshal(data, &tmp); err == nil {
				return tmp
			}
			return nil
		default:
			// Fallback: return as-is (will be JSON marshaled later)
			return meta
	}
}
