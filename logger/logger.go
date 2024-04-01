package logger

import (
	"context"
	"dataset"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/**
logger has levels: Panic, Fatal, Error, Warn, Info, and Debug
Panic will log a message and then panic.  This is not to be used in production.
Fatal will log a message and Goexit(), but osExit() if the caller is context.Background()
Fatal should also be used in rare cases.
Error will log a message an return, but it is expected that the transaction will fail.
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

var logLevel LogLevel = LOGDEBUG
var dumpSkipLines = 3 // number of dump lines to skip

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

func SetDumpSkipLines(lines int) {
	dumpSkipLines = lines
}

// Panic will log the message with Println and then call panic()
func Panic(ctx context.Context, param ...any) {
	panicLog.Panicln(param, requestInfo(ctx))
}

// Fatal will log the message with Println and call os.Exit is in background
func Fatal(ctx context.Context, param ...any) {
	fatalLog.Println(param, requestInfo(ctx))
	fatalLog.Println(dumpLines())
	if ctx == context.Background() {
		os.Exit(1)
	} else {
		runtime.Goexit()
	}
}

// Error will log the message with Println
func Error(ctx context.Context, param ...any) {
	errorLog.Println(param, requestInfo(ctx))
	errorLog.Println(dumpLines())
}

// Warn will log the message with Println and then continue
func Warn(ctx context.Context, param ...any) {
	if logLevel >= LOGWARN {
		warnLog.Println(param, requestInfo(ctx))
		warnLog.Println(dumpLines())
	}
}

// Info will log the message with Println and then continue
func Info(ctx context.Context, param ...any) {
	if logLevel >= LOGINFO {
		infoLog.Println(param, requestInfo(ctx))
	}
}

// Debug will log the message with Println and then continue
func Debug(ctx context.Context, param ...any) {
	if logLevel >= LOGDEBUG {
		debugLog.Println(param, requestInfo(ctx))
	}
}

func requestInfo(ctx context.Context) string {
	if ctx != nil {
		request := ctx.Value("request")
		if request != nil {
			req := request.(dataset.RequestType)
			result := `AudioSource=` + string(req.AudioSource) + ` TextSource=` + string(req.TextSource) +
				` Testament=` + string(req.Testament)
			return result
		}
	}
	return ""
}

func dumpLines() string {
	var results = make([]string, 0)
	var pcs = make([]uintptr, 7, 7)
	num := runtime.Callers(dumpSkipLines, pcs)
	for i := 0; i < num; i++ {
		level := strconv.Itoa(dumpSkipLines + i)
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
