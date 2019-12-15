/*
* @Author:	Payton
* @Date:	buffer_pool
* @DESC: 	DESC
 */

package gen_log

import (
	"sync"
)

type LogBufferPool struct {
	pool_ *sync.Pool
}

func new_pool() *LogBufferPool {
	return &LogBufferPool{pool_: &sync.Pool{
		New: func() interface{} {
			return &LogBuffer{buff_: make([]byte, 0, size_)}
		},
	}}
}

var (
	logbuffer_pool *LogBufferPool = new_pool()
)

func (buffer_pool *LogBufferPool) Pop() *LogBuffer {
	buff_ := buffer_pool.pool_.Get().(*LogBuffer)
	buff_.pool = buffer_pool
	buff_.Reset()
	return buff_
}

func (buffer_pool *LogBufferPool) Push(buffer *LogBuffer) {
	buffer_pool.pool_.Put(buffer)
}
