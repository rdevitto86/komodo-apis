package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"komodo-forge-apis-go/config"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

type LoggerConfig struct {
	AppName          string
	LogLevel         string
}

var (
	initOnce sync.Once
	logLevel slog.LevelVar
	appName	string
	slogger *slog.Logger
)

// Initialize the logger with the specified configuration.
func Init(name string, level string) {
	initOnce.Do(func() { slogger = initialize(name, level) })
}

func Info(msg string, details ...any)  { writeLog("INFO", msg, details...) }
func Warn(msg string, details ...any)  { writeLog("WARN", msg, details...) }
func Error(msg string, details ...any) { writeLog("ERROR", msg, details...) }
func Fatal(msg string, details ...any) { writeLog("FATAL", msg, details...) }
func Debug(msg string, details ...any) { writeLog("DEBUG", msg, details...) }
func Trace(msg string, details ...any) { writeLog("TRACE", msg, details...) }

// initialize sets up the global logger with the specified configuration.
func initialize(name string, lvl string) *slog.Logger {
	appName = name
	if name == "" {
		appName = config.GetConfigValue("APP_NAME")
	}

	if lvl != "" {
		logLevel.Set(parseSlogLevel(lvl))
	} else if lvl := config.GetConfigValue("LOG_LEVEL"); lvl != "" {
		logLevel.Set(parseSlogLevel(lvl))
	} else if strings.ToLower(config.GetConfigValue("ENV")) == "local" {
		logLevel.Set(parseSlogLevel("DEBUG"))
	} else {
		logLevel.Set(parseSlogLevel("ERROR"))
	}

	slogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &logLevel}))
	slogger.Info("runtime logger initialized", slog.String("app", appName))
	return slogger
}

// Helper function to write log entries with level filtering.
func writeLog(level, msg string, details ...any) {
	initOnce.Do(func() {
		initialize(config.GetConfigValue("APP_NAME"), config.GetConfigValue("LOG_LEVEL"))
	})

	lvl := parseSlogLevel(level)

	// Filter by log level. Skip if below global threshold.
	if lvl < logLevel.Level() { return }

	if len(details) == 0 {
		slogger.LogAttrs(context.Background(), lvl, msg, slog.String("app", appName))
		return
	}

	// Handle all details as key-value pairs
	attrs := make([]slog.Attr, 0, len(details)/2)
	attrs = append(attrs, slog.String("app", appName))

	for i := 0; i < len(details); i += 2 {
		if i+1 >= len(details) {
			// Handle case where we have an odd number of arguments
			attrs = append(attrs, slog.Any(fmt.Sprintf("arg%d", i), normalizeMeta(details[i])))
			break
		}
		
		key, ok := details[i].(string)
		if !ok {
			// If key is not a string, use a generic key
			key = fmt.Sprintf("arg%d", i)
			i-- // Adjust index since we didn't consume the value
			continue
		}
		
		value := normalizeMeta(details[i+1])
		attrs = append(attrs, slog.Any(key, value))
	}
	slogger.LogAttrs(context.Background(), lvl, msg, attrs...)
}

// Handles various common types and converts them to a format suitable for slog.
func normalizeMeta(v any) any {
	if v == nil { return nil }
	
	// Handle common types explicitly
	switch val := v.(type) {
		case string, bool, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, float32, float64:
			return val
		
		case error:
			return val.Error()
				
    case time.Time:
      return val.UTC().Format(time.RFC3339Nano)
            
    case *time.Time:
			if val != nil {
				return val.UTC().Format(time.RFC3339Nano)
			}
			return nil
            
    case fmt.Stringer:
      return val.String()
				
		case *http.Request:
			return map[string]any{
				"method":      val.Method,
				"path":        val.URL.Path,
				"host":        val.Host,
				"request_id":  val.Context().Value("request_id"),
				"pathParams":  val.Context().Value("pathParams"),
				"queryParams": val.URL.Query(),
				"headers":     val.Header,
			}
				
		case json.RawMessage:
			var tmp any
			if err := json.Unmarshal(val, &tmp); err == nil {
				return normalizeMeta(tmp)
			}
			return string(val)
				
		case []byte:
			var tmp any
			if err := json.Unmarshal(val, &tmp); err == nil {
				return normalizeMeta(tmp)
			}
			return string(val)
	}
	
	// Handle maps, slices, and arrays
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
		case reflect.Map:
			result := make(map[string]any, rv.Len())
			for _, key := range rv.MapKeys() {
				if strKey, ok := key.Interface().(string); ok {
					result[strKey] = normalizeMeta(rv.MapIndex(key).Interface())
				}
			}
			return result
			
		case reflect.Slice, reflect.Array:
			result := make([]any, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				result[i] = normalizeMeta(rv.Index(i).Interface())
			}
			return result
				
		case reflect.Ptr:
			if rv.IsNil() { return nil }
			return normalizeMeta(rv.Elem().Interface())
				
		case reflect.Struct:
			if rv.CanInterface() {
				// Handle common structs that implement Stringer or have custom formatting
				if s, ok := v.(fmt.Stringer); ok {
					return s.String()
				}
					
				// For other structs, convert to map
				result := make(map[string]any)
				t := rv.Type()

				for i := 0; i < t.NumField(); i++ {
					field := t.Field(i)
					// Skip unexported fields
					if field.PkgPath != "" { continue }
					
					// Handle json tags if present
					name := field.Name
					if tag := field.Tag.Get("json"); tag != "" {
						if tag == "-" { continue }
						if comma := strings.Index(tag, ","); comma != -1 {
							name = tag[:comma]
						} else {
							name = tag
						}
					}
					
					fieldVal := rv.Field(i)
					if fieldVal.CanInterface() {
						result[name] = normalizeMeta(fieldVal.Interface())
					}
				}
				return result
			}
	}
	return fmt.Sprintf("%+v", v)
}

// Converts a string level to slog.Level
func parseSlogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
		case "TRACE":
			return slog.Level(-8)
		case "DEBUG":
			return slog.LevelDebug
		case "INFO":
			return slog.LevelInfo
		case "WARN":
			return slog.LevelWarn
		case "ERROR":
			return slog.LevelError
		case "FATAL":
			return slog.LevelError
		default:
			return slog.LevelError
	}
}
