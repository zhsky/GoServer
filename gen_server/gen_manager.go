/*
* @Author:	Payton
* @Date:	gen_manager
* @DESC: 	DESC
 */

package gen_server

import (
	"errors"
	"sync"
	"time"
)

var (
	server_map = make(map[string]*Channel)
	sync_wait  sync.WaitGroup
	sm_lock    sync.Mutex
	close      bool = false
)

var (
	CLOSING = errors.New("CLOSING")
)

func regist_server(server Server, channel *Channel) error {
	sm_lock.Lock()
	defer sm_lock.Unlock()
	if close {
		return errors.New("CLOSING")
	}
	var ServerName = server.Module()
	if _, Exist := server_map[ServerName]; Exist {
		return errors.New("Repeat Server " + ServerName)
	}

	server_map[ServerName] = channel
	sync_wait.Add(1)
	return nil
}

func unregist_server(server Server) {
	var ServerName = server.Module()
	sm_lock.Lock()
	delete(server_map, ServerName)
	sm_lock.Unlock()
	sync_wait.Done()
}

func Open() {
	close = false
}

func SyncClose() {
	close = true
	sm_lock.Lock()
	var async_msg = AsyncMsg{}
	for _, channel := range server_map {
		async_msg.MsgId = -1 //===Stop====
		channel.MessqgeQueue <- async_msg
	}
	sm_lock.Unlock()
	sync_wait.Wait()
}

func SendAsyncMsg(ServerName string, msg interface{}) error {
	if close {
		return CLOSING
	}
	if channel, Exist := server_map[ServerName]; Exist {
		channel.MessqgeQueue <- msg
	} else {
		return errors.New("Server Not Exist " + ServerName)
	}
	return nil
}

func SendSyncMsg(ServerName string, msg interface{}) (reply []interface{}, err error) {
	return SendSyncMsgTimeout(ServerName, msg, 5000)
}

func SendSyncMsgTimeout(ServerName string, msg interface{}, timeout time.Duration) (reply []interface{}, err error) {
	if close {
		return nil, CLOSING
	}
	if channel, Exist := server_map[ServerName]; Exist {
		var reply_channel = make(chan []interface{}, 1)
		var sync_msg = SyncMsg{
			AsyncMsg:         msg.(AsyncMsg),
			callback_channel: reply_channel,
		}
		channel.SyncMessqgeQueue <- sync_msg

		select {
		case reply = <-reply_channel:
		case <-time.After(time.Millisecond * timeout):
			return nil, errors.New("SYNC timeout")
		}
	} else {
		return nil, errors.New("Server Not Exist " + ServerName)
	}
	return reply, nil
}

func Call(ServerName string, fun CallbackFunc) error {
	if close {
		return CLOSING
	}
	if channel, Exist := server_map[ServerName]; Exist {
		channel.callfunc(fun)
	} else {
		return errors.New("Server Not Exist " + ServerName)
	}

	return nil
}
