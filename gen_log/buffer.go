/*
* @Author:	Payton
* @Date:	buffer
* @DESC: 	DESC
 */

package gen_log

import (
	"reflect"
)

var (
	size_ int32 = 1024
)

type LogBuffer struct {
	buff_ []byte
	pool  *LogBufferPool
}

func (buffer *LogBuffer) Write() {
	write_log(buffer)
}

func (buffer *LogBuffer) Reset() {
	buffer.buff_ = buffer.buff_[:0]
}

func (buffer *LogBuffer) String() string {
	return string(buffer.buff_)
}

func (buffer *LogBuffer) EnsureCanRead(size int) {
	if len(buffer.buff_) < size {
		b := make([]byte, size-len(buffer.buff_), size-len(buffer.buff_))
		buffer.buff_ = append(buffer.buff_, b...)
	}
}

func (buffer *LogBuffer) Log(args ...interface{}) *LogBuffer {
	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			buffer = buffer.AppendInt(int64(v)).AppendByte(' ')
		case int32:
			buffer = buffer.AppendInt(int64(v)).AppendByte(' ')
		case int64:
			buffer = buffer.AppendInt(v).AppendByte(' ')
		case uint:
			buffer = buffer.AppendUInt(uint64(v)).AppendByte(' ')
		case uint32:
			buffer = buffer.AppendUInt(uint64(v)).AppendByte(' ')
		case uint64:
			buffer = buffer.AppendUInt(v).AppendByte(' ')
		case byte:
			buffer = buffer.AppendByte(v).AppendByte(' ')
		case bool:
			buffer = buffer.AppendBool(v).AppendByte(' ')
		case string:
			buffer = buffer.AppendString(v).AppendByte(' ')
		default:
			buffer = buffer.AppendString("!!!ERROR TYPE:" + reflect.TypeOf(arg).String()).AppendByte(' ')
		}
	}
	return buffer
}

func (buffer *LogBuffer) AppendByte(v byte) *LogBuffer {
	buffer.buff_ = append(buffer.buff_, v)
	return buffer
}

func (buffer *LogBuffer) AppendString(str string) *LogBuffer {
	buffer.buff_ = append(buffer.buff_, str...)
	return buffer
}

func (buffer *LogBuffer) AppendUInt(v uint64) *LogBuffer {
	buffer.buff_ = appendInt(buffer.buff_, v)
	return buffer
}

func (buffer *LogBuffer) AppendInt(v int64) *LogBuffer {
	if v < 0 {
		buffer.buff_ = append(buffer.buff_, '-')
		v = -v
	}
	return buffer.AppendUInt(uint64(v))
}

func (buffer *LogBuffer) AppendBool(v bool) *LogBuffer {
	if v {
		buffer.buff_ = append(buffer.buff_, "true"...)
	} else {
		buffer.buff_ = append(buffer.buff_, "false"...)
	}
	return buffer
}

func appendInt(buff_ []byte, x uint64) []byte {
	var a [20]byte
	i := 19
	y, z := x/10, x%10
	for y > 9 {
		a[i] = byte('0' + z)
		i--
		x = y
		y, z = x/10, x%10
	}

	if z > 0 {
		a[i] = byte('0' + z)
		i--
	}
	if y > 0 {
		a[i] = byte('0' + y)
		i--
	}
	buff_ = append(buff_, a[i+1:]...)
	return buff_
}
