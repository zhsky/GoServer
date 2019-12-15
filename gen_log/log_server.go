/*
* @Author:	Payton
* @Date:	log_server
* @DESC: 	DESC
 */

package gen_log

import (
	"io"
	"log"
	"os"
	"time"
)

const (
	FILE_SZIE = 1024 * 1024 * 500
)

var (
	recv_chan      = make(chan *LogBuffer, 10000)
	size      int  = 0
	stop      bool = true
	file      *os.File

	stop_chan           = make(chan bool)
	filename_pre string = "./logs/log_"
)

func StartLog(level int) {
	SetRunLevel(level)
	stop = false
	go func() {
		recv_chan <- logbuffer_pool.Pop().AppendString("LOG START!\n")
		for !stop || len(recv_chan) > 0 {
			run()
		}
		close_file()
		stop_chan <- true
	}()
}

func SetRunLevel(level int) {
	run_level = level
}

func StopLog() {
	stop = true
	recv_chan <- logbuffer_pool.Pop().AppendString("NORMAL STOP!\n")
	<-stop_chan
}

func write_log(buffer *LogBuffer) {
	if !stop {
		recv_chan <- buffer
	}
}
func run() {
	defer func() {
		err := recover()
		if err != nil {
			log.Print(err)
		}
	}()
	if file == nil {
		var err error
		file, err = os.OpenFile(get_file_name(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Print(err)
			return
		}
	}

	for !stop || len(recv_chan) > 0 {
		select {
		case str := <-recv_chan:
			s := str.buff_[:]
			n, err := file.Write(s)
			size += n
			for err == io.ErrShortWrite {
				s = s[n:]
				n, err = file.Write(s)
				size += n
			}
			logbuffer_pool.Push(str)
			str = nil

			if err != nil {
				log.Print(err)
				close_file()
				return
			}

			if size > FILE_SZIE {
				close_file()
				return
			}
		}
	}
	return
}

func close_file() {
	if file != nil {
		file := file
		go func() {
			file.Sync()
			file.Close()
		}()
	}
	file = nil
	size = 0
	return
}

func get_file_name() string {
	ftime := time.Now().Format("20060102-15.0405") //"2006-01-02T15:04:05Z07:00"
	return filename_pre + ftime + ".log"
}
