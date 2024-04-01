package logger

import (
	"context"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/**
logger has levels: Panic, Fatal, Error, Warn, Info, and Debug
Panic will log a message and then panic.  Fatal will log a message and then exit.
Error will Goexit() when a goroutine calls the logger, and os.Exit() when the
main goroutine calls it.  It determines if a routine is main by checking for
background context.
Warn, Info, and Debug log messages and continue.
*/

type LogLevel int

const (
	LOGFATAL LogLevel = iota + 1
	LOGERROR
	LOGWARN
	LOGINFO
	LOGDEBUG
)

var logLevel LogLevel = LOGINFO

var panicLog *log.Logger
var fatalLog *log.Logger
var errorLog *log.Logger
var warnLog *log.Logger
var infoLog *log.Logger
var debugLog *log.Logger

func init() {
	setFile(os.Stderr)
}

// SetOutput accepts: stdout, stderr, or a filePath
func SetOutput(ctx context.Context, filePath string) {
	if filePath == "stdout" {
		setFile(os.Stdout)
	} else if filePath == "stderr" {
		setFile(os.Stderr)
	} else {
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			Warn(ctx, 0, "log.SetOutput failed", err)
		} else {
			setFile(file)
		}
	}
}

func setFile(file *os.File) {
	panicLog = log.New(file, "PANIC ", log.Ldate|log.Ltime)
	fatalLog = log.New(file, "FATAL ", log.Ldate|log.Ltime)
	errorLog = log.New(file, "ERROR ", log.Ldate|log.Ltime)
	warnLog = log.New(file, "WARN ", log.Ldate|log.Ltime)
	infoLog = log.New(file, "INFO ", log.Ldate|log.Ltime)
	debugLog = log.New(file, "DEBUG ", log.Ldate|log.Ltime|log.Lmicroseconds)
}

// SetLevel set an error reporting level. The logger will process messages of this level and higher.
func SetLevel(level LogLevel) {
	logLevel = level
}

// Panic will log the message with Println and then call panic()
func Panic(ctx context.Context, skip int, param ...any) {
	panicLog.Panicln(fileLine(skip), param, requestInfo(ctx))
}

// Fatal will log the message with Println and then call os.Exit(1)
func Fatal(ctx context.Context, skip int, param ...any) {
	fatalLog.Fatalln(fileLine(skip), param, requestInfo(ctx))
}

// Error will log the message with Println and then call runtime.Goexit()
func Error(ctx context.Context, skip int, param ...any) {
	if logLevel >= LOGERROR {
		errorLog.Println(fileLine(skip), param, requestInfo(ctx))
	}
	if ctx == context.Background() {
		os.Exit(1)
	} else {
		runtime.Goexit()
	}
}

// Warn will log the message with Println and then continue
func Warn(ctx context.Context, skip int, param ...any) {
	if logLevel >= LOGWARN {
		warnLog.Println(param, requestInfo(ctx))
	}
}

// Info will log the message with Println and then continue
func Info(ctx context.Context, skip int, param ...any) {
	if logLevel >= LOGINFO {
		infoLog.Println(param, requestInfo(ctx))
	}
}

// Debug will log the message with Println and then continue
func Debug(ctx context.Context, skip int, param ...any) {
	if logLevel >= LOGDEBUG {
		debugLog.Println(param, requestInfo(ctx))
	}
}

func requestInfo(ctx context.Context) string {
	if ctx != nil {
		request := ctx.Value("request")
		if request != nil {
			return "Request:" + request.(string)
		}
	}
	return ""
}

func fileLine(skip int) string {
	start := 3 //skip + 2
	var results = make([]string, 0)
	var pcs = make([]uintptr, 7, 7)
	num := runtime.Callers(start, pcs)
	for i := 0; i < num; i++ {
		level := strconv.Itoa(start + i)
		fun := runtime.FuncForPC(pcs[i])
		file, line := fun.FileLine(pcs[i])
		start := strings.LastIndex(file, "/") + 1
		base := file[start:]
		fileLine := base + ":" + strconv.Itoa(line)
		name := fun.Name()
		start = strings.LastIndex(name, "/") + 1
		fname := name[start:]
		results = append(results, level+" "+fileLine+" "+fname)
	}
	return strings.Join(results, "\n")
}
