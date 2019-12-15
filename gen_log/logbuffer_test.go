/*
* @Author:	Payton
* @Date:	logbuffer_test
* @DESC: 	DESC
 */

package gen_log

import (
	"testing"
)

func TestLogBuffer(t *testing.T) {
	StartLog(LOG_DEBUG)

	Info("TestLogBuffer ", 2, "TestLogBuffer ", '9')
	Error(9223372036854775807, " TestLogBuffer ", int16(2), true, '9')
	Debug(int64(-9223372036854775807), " TestLogBuffer ", int64(2), true, '9')
	Fatal(uint64(18446744073709551615), " TestLogBuffer ", int64(2), true, '9')

	StopLog()
}

func BenchmarkBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info(uint64(18446744073709551615), " TestLogBuffer ", int64(2), true, '9')
	}
}
