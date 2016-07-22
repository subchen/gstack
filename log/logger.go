package log

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Log Level
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	levelStr = []string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
		"FATAL",
	}
	levelStrWithColor = []string{
		"\033[34mDEBUG\033[0m",
		"\033[32mINFO\033[0m",
		"\033[33mWARN\033[0m",
		"\033[31mERROR\033[0m",
		"\033[35mFATAL\033[0m",
	}
)

func New(out io.Writer) *Logger {
	return &Logger{
		out:        out,
		level:      DEBUG,
		name:       "",
		timeLayout: "2006-01-02 15:04:05.999",
		longFile:   false,
		colorful:   false,
		callerSkip: 2,
	}
}

type Logger struct {
	mu         sync.Mutex
	out        io.Writer
	level      int
	name       string
	timeLayout string
	longFile   bool
	colorful   bool
	callerSkip int
}

func (l *Logger) GetLevel() int {
	return l.level
}

func (l *Logger) SetLevel(level int) *Logger {
	l.level = level
	return l
}

func (l *Logger) SetName(name string) *Logger {
	l.name = name
	return l
}

func (l *Logger) SetTimeLayout(layout string) *Logger {
	l.timeLayout = layout
	return l
}

func (l *Logger) UseLongFile() *Logger {
	l.longFile = true
	return l
}

func (l *Logger) SetColorful(colorful bool) *Logger {
	l.colorful = colorful
	return l
}

func (l *Logger) SkipCaller(skip int) *Logger {
	l.callerSkip = skip
	return l
}

func (l *Logger) SetWriter(w io.Writer) *Logger {
	l.out = w
	return l
}

func (l *Logger) IsDebugEnabled() bool {
	return l.level <= DEBUG
}

func (l *Logger) IsInfoEnabled() bool {
	return l.level <= INFO
}

func (l *Logger) IsWarnEnabled() bool {
	return l.level <= WARN
}

func (l *Logger) IsErrorEnabled() bool {
	return l.level <= ERROR
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= DEBUG {
		l.log(DEBUG, msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, msg, args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= WARN {
		l.log(WARN, msg, args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, msg, args...)
	}
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(FATAL, msg, args...)
}

func (l *Logger) log(level int, msg string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(l.callerSkip)
	if !ok {
		file = "???"
		line = 0
	} else if !l.longFile {
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	}

	buf := new(bytes.Buffer)
	buf.WriteByte(' ')
	if l.name != "" {
		buf.WriteString(l.name)
		buf.WriteByte(' ')
	}
	if l.colorful {
		buf.WriteString(levelStrWithColor[level])
	} else {
		buf.WriteString(levelStr[level])
	}
	buf.WriteByte(' ')
	buf.WriteString(file)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(line))
	buf.WriteByte(' ')
	fmt.Fprintf(buf, msg, args...)
	buf.WriteByte('\n')

	if level == FATAL {
		for i := l.callerSkip; ; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			fmt.Fprintf(buf, "\tat %s:%d (0x%x)\n", file, line, pc)
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timeStr := time.Now().Format(l.timeLayout)

	l.out.Write([]byte(timeStr))
	l.out.Write(buf.Bytes())
}
