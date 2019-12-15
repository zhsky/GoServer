/*
* @Author:	Payton
* @Date:	gen_server
* @DESC: 	DESC
 */

package gen_server

import (
	"gen_log"
	"sync"
	"time"
	"util"
)

type Server interface {
	Module() string
	Init()
	HandleMsg(msg interface{}) int32
	HandleSyncMsg(msg interface{}) []interface{} //传值返回
	HandleTimeout(int64) int32
	Stop(exit_code int32)
}

type AsyncMsg struct {
	MsgId int32
	Args  []interface{}
}

type SyncMsg struct {
	AsyncMsg
	callback_channel chan []interface{}
}
type Channel struct {
	MessqgeQueue     chan interface{}
	SyncMessqgeQueue chan interface{}
	timeout          <-chan time.Time
	callbacks        []util.CallbackFunc
	cf_lock          sync.Mutex
}

func Start(server Server) {
	go util.RunFuncWithPainc(func() {
		gen_init(server)
	})
}

func gen_init(server Server) {
	server.Init()
	var MsgSize = 10000
	var channel = &Channel{
		MessqgeQueue:     make(chan interface{}, MsgSize),
		SyncMessqgeQueue: make(chan interface{}, MsgSize),
		callbacks:        make([]CallbackFunc, 0, MsgSize),
		timeout:          time.Tick(5000),
	}

	if err := regist_server(server, channel); err == nil {
		var ret int32 = 0
		for ret == 0 {
			ret = channel.enter_loop(server)
		}
		unregist_server(server)
		server.Stop(ret)
	} else {
		gen_log.Error(err)
	}
}

func (channel *Channel) enter_loop(server Server) (retcode int32) {
	defer func() {
		err := recover()
		if err != nil {
			gen_log.Error(err)
		}
	}()

	var ret int32 = 0
	for {
		select {
		case msg := <-channel.MessqgeQueue:
			if ret = server.HandleMsg(msg); ret != 0 {
				return ret
			}
		case msg := <-channel.SyncMessqgeQueue:
			handle_syncsg(server, msg)
		case now := <-channel.timeout:
			if ret = server.HandleTimeout(now.Unix()); ret != 0 {
				return ret
			}
			channel.handle_callfunc()
		}
	}
}

func handle_syncsg(server Server, msg interface{}) {
	var sync_msg = msg.(SyncMsg)
	reply := server.HandleSyncMsg(sync_msg.AsyncMsg)
	sync_msg.callback_channel <- reply
}

func (channel *Channel) callfunc(fun util.CallbackFunc) {
	channel.cf_lock.Lock()
	channel.callbacks = append(channel.callbacks, fun)
	channel.cf_lock.Unlock()
}

func (channel *Channel) handle_callfunc() {
	if len(channel.callbacks) > 0 {
		channel.cf_lock.Lock()
		callbacksCopy := channel.callbacks
		channel.callbacks = make([]util.CallbackFunc, 0, len(channel.callbacks))
		channel.cf_lock.Unlock()

		for _, Fun := range callbacksCopy {
			util.RunFuncWithPainc(Fun)
		}
	}
}
