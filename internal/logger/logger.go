package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"gogurt/internal/types"
)

// Level is an alias for the Level type in the types package.
// We also expose constants here so other packages can use logger.INFO, etc.
type Level = types.LogLevel

const (
	INFO    = types.INFO
	WARNING = types.WARNING
	ERROR   = types.ERROR
	FATAL   = types.FATAL
)

// RuntimeInfo holds details about the application's runtime environment.
type RuntimeInfo struct {
	GoVersion string `json:"go"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// LogEntry represents a structured log entry for the log file.
type LogEntry struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Runtime   RuntimeInfo    `json:"runtime"`
	Caller    string         `json:"caller,omitempty"`
	Message   string         `json:"message"`
	Fields    map[string]any `json:"fields,omitempty"`
	TraceID   string         `json:"trace_id,omitempty"`
}

// String provides a human-readable, plain-text format for the console.
func (e LogEntry) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] [%s] [%s] %s", e.Timestamp, e.Level, e.Caller, e.Message))

	if len(e.Fields) > 0 {
		sb.WriteString(" | Fields:")
		for k, v := range e.Fields {
			sb.WriteString(fmt.Sprintf(" %s=%v", k, v))
		}
	}
	return sb.String()
}

// Logger provides structured logging capabilities with distinct outputs.
type Logger struct {
	consoleWriter io.Writer
	fileWriter    io.Writer
	consoleFormat types.LogFormat
	fileFormat    types.LogFormat
}

var defaultLogger *Logger

// Init initializes the default logger with a metrics service and format configs.
func Init(consoleFormat, fileFormat types.LogFormat) {
	logDir := "logs"
	logPath := logDir + "/godash.log"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	defaultLogger = &Logger{
		consoleWriter: os.Stdout,
		fileWriter:    logFile,
		consoleFormat: consoleFormat,
		fileFormat:    fileFormat,
	}
}

// getCallerInfo inspects the call stack to find the original caller of the log function.
func getCallerInfo() string {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return "???"
	}

	fn := runtime.FuncForPC(pc)
	funcName := "???"
	if fn != nil {
		funcName = path.Base(fn.Name())
	}

	modulePath := file
	if idx := strings.LastIndex(file, "godash/"); idx != -1 {
		modulePath = file[idx:]
	}

	return fmt.Sprintf("%s:%d %s", modulePath, line, funcName)
}

// Log logs a message with the specified level.
func Log(level Level, message string) {
	logStructured(defaultLogger, level, message, nil, "", getCallerInfo())
}

// LogWithContext logs a message with context and additional fields.
func LogWithContext(ctx context.Context, level Level, message string, fields map[string]any) {
	traceID := extractTraceID(ctx)
	logStructured(defaultLogger, level, message, fields, traceID, getCallerInfo())
}

// logStructured constructs the log entry and writes it to the console and file.
func logStructured(l *Logger, level Level, message string, fields map[string]any, traceID string, caller string) {
	if level >= ERROR {
		log.Println(message)
	}
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level.String(),
		Runtime: RuntimeInfo{
			GoVersion: runtime.Version(),
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		},
		Caller:  caller,
		Message: message,
		Fields:  fields,
		TraceID: traceID,
	}

	if l.consoleFormat == types.FormatJSON {
		jsonData, err := json.Marshal(entry)
		if err == nil {
			fmt.Fprintln(l.consoleWriter, string(jsonData))
		}
	} else {
		fmt.Fprintln(l.consoleWriter, entry.String())
	}

	if l.fileFormat == types.FormatJSON {
		jsonData, err := json.Marshal(entry)
		if err == nil {
			fmt.Fprintln(l.fileWriter, string(jsonData))
		}
	} else {
		fmt.Fprintln(l.fileWriter, entry.String())
	}

	if level == FATAL {
		os.Exit(1)
	}
}

// LogError is a helper for logging errors.
func LogError(ctx context.Context, message string, err error, fields map[string]any) {
	if err == nil {
		return
	}
	if fields == nil {
		fields = make(map[string]any)
	}
	fields["error"] = err.Error()
	LogWithContext(ctx, ERROR, message, fields)
}

// extractTraceID gets a trace ID from context, if present.
func extractTraceID(ctx context.Context) string {
	const traceKey = "trace_id"
	if ctx == nil {
		return ""
	}
	val := ctx.Value(traceKey)
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func NewLogger(consoleWriter io.Writer, fileWriter io.Writer, consoleFormat types.LogFormat, fileFormat types.LogFormat) *Logger {
	return &Logger{
		consoleWriter: consoleWriter,
		fileWriter:    fileWriter,
		consoleFormat: consoleFormat,
		fileFormat:    fileFormat,
	}
}

var (
	loggerInstance *Logger
	once           sync.Once
)

// GetLogger returns the singleton Logger instance.
func GetLogger() *Logger {
	once.Do(func() {
		loggerInstance = NewLogger(os.Stdout, os.Stdout, types.FormatText, types.FormatText)
	})
	return loggerInstance
}

