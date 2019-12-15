/*
* @Author:	Payton
* @Date:	log_wrapper
* @DESC: 	DESC
 */

package gen_log

import (
	"runtime"
	"time"
)

const (
	LOG_INFO  = 1
	LOG_FATAL = 2
	LOG_DEBUG = 3

	FORMAT = `15:04:05`
	INFO   = `[INFO] `
	ERROR  = `[ERROR] `
	DEBUG  = `[DEBUG] `
	FATAL  = `FATAL] `
)

var (
	run_level int = LOG_FATAL
)

func Info(args ...interface{}) {
	logbuffer_pool.Pop().AppendString(time.Now().Format(FORMAT)).
		AppendString(INFO).AppendByte('\t').
		Log(args...).AppendByte('\n').Write()
}

func Error(args ...interface{}) {
	_, name, line, _ := runtime.Caller(1)
	logbuffer_pool.Pop().AppendString(time.Now().Format(FORMAT)).
		AppendString(ERROR).AppendString(name).AppendByte(':').AppendInt(int64(line)).AppendByte('\t').
		Log(args...).AppendByte('\n').Write()
}

func Debug(args ...interface{}) {
	if run_level < LOG_DEBUG {
		return
	}
	_, name, line, _ := runtime.Caller(1)
	logbuffer_pool.Pop().AppendString(time.Now().Format(FORMAT)).
		AppendString(DEBUG).AppendString(name).AppendByte(':').AppendInt(int64(line)).AppendByte('\t').
		Log(args...).AppendByte('\n').Write()
}

func Fatal(args ...interface{}) {
	if run_level < LOG_FATAL {
		return
	}
	_, name, line, _ := runtime.Caller(1)
	logbuffer_pool.Pop().AppendString(time.Now().Format(FORMAT)).
		AppendString(FATAL).AppendString(name).AppendByte(':').AppendInt(int64(line)).AppendByte('\t').
		Log(args...).AppendByte('\n').Write()
	Trace()
}

func Trace() {
	var log = logbuffer_pool.Pop().AppendString(time.Now().Format(FORMAT))
	log.EnsureCanRead(1024)
	n := runtime.Stack(log.buff_, false)
	log.buff_ = log.buff_[0:n]
	log.AppendByte('\n').Write()
}
