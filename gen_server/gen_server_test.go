/*
* @Author:	Payton
* @Date:	gen_server_test
* @DESC: 	DESC
 */

package gen_server

import (
	"fmt"
	"testing"
	"time"
)

type TestModule struct {
}

func (server *TestModule) Module() string {
	return "TestModule"
}

func (server *TestModule) Init() {
	fmt.Println("TestModule Init")
}
func (server *TestModule) Stop(exit_code int32) {
	fmt.Println("TestModule Stop", exit_code)
}

func (server *TestModule) HandleMsg(msg interface{}) (ret int32) {
	ret = 0
	async_msg, err := msg.(AsyncMsg)
	if !err {
		return
	}

	switch async_msg.MsgId {
	case -1:
		ret = 1
	case 1001:
		handle_1001(async_msg.Args)
	case 1002:
		handle_1002(async_msg.Args)
	default:
		fmt.Println("ERROR MsgId ", async_msg.MsgId)
	}
	return
}

func (server *TestModule) HandleSyncMsg(msg interface{}) (reply []interface{}) {
	sync_msg, err := msg.(AsyncMsg)
	if !err {
		return
	}

	switch sync_msg.MsgId {
	case 1002:
		reply = handle_1002(sync_msg.Args)
	case 1003:
		reply = handle_1003(sync_msg.Args)
	default:
		fmt.Println("ERROR MsgId ", sync_msg.MsgId)
		reply = nil
	}
	return
}

func (server *TestModule) HandleTimeout(now int64) (ret int32) {
	ret = 0
	if now%1000000 == 0 {
		fmt.Println("ERROR MsgId ", now)
	}
	return
}

// args[0] int
// args[1] string
func handle_1001(args []interface{}) {
	arg0, err0 := args[0].(int)
	arg1, err1 := args[1].(string)
	if err0 && err1 {
		fmt.Println("handle_1001 ", arg0, arg1)
	}
}

// args[0] string
func handle_1002(args []interface{}) []interface{} {
	arg0, err0 := args[0].(string)
	if err0 {
		fmt.Println("handle_1002 ", arg0)
	}
	var reply = make([]interface{}, 0, 2)
	reply = append(reply, 100)
	reply = append(reply, "bobo")
	return reply
}

// args[0] int
// args[1] string
// args[2] int
func handle_1003(args []interface{}) []interface{} {
	arg0, err0 := args[0].(int)
	arg1, err1 := args[1].(string)
	arg2, err2 := args[2].(int)
	if err0 && err1 && err2 {
		fmt.Println("handle_1003 ", arg0, arg1, arg2)
	}

	var reply = make([]interface{}, 0, 3)
	reply = append(reply, 100)
	reply = append(reply, "bobo")
	reply = append(reply, 888)
	return reply
}

func test_send_msg() {
	var async_msg = AsyncMsg{}

	//===Async msg====
	async_msg.MsgId = 1001
	async_msg.Args = make([]interface{}, 0, 2)
	async_msg.Args = append(async_msg.Args, 100)
	async_msg.Args = append(async_msg.Args, "Async msg 1001")
	SendAsyncMsg("TestModule", async_msg)

	//===Async msg====
	async_msg.MsgId = 1002
	async_msg.Args = make([]interface{}, 0, 1)
	async_msg.Args = append(async_msg.Args, "Async msg 1002")
	SendAsyncMsg("TestModule", async_msg)

	//===Sync msg====
	async_msg.MsgId = 1002
	async_msg.Args = make([]interface{}, 0, 1)
	async_msg.Args = append(async_msg.Args, "Sync msg 1002")
	if reply, err := SendSyncMsg("TestModule", async_msg); err == nil {
		arg0, err0 := reply[0].(int)
		arg1, err1 := reply[1].(string)
		if err0 && err1 {
			fmt.Println("1002 reply ", arg0, arg1)
		}
	} else {
		fmt.Println("1002 reply ", err)
	}

	//===Sync msg====
	async_msg.MsgId = 1003
	async_msg.Args = make([]interface{}, 0, 2)
	async_msg.Args = append(async_msg.Args, 100)
	async_msg.Args = append(async_msg.Args, "Sync msg 1003")
	async_msg.Args = append(async_msg.Args, 999)
	if reply, err := SendSyncMsg("TestModule", async_msg); err == nil {
		arg0, err0 := reply[0].(int)
		arg1, err1 := reply[1].(string)
		arg2, err2 := reply[2].(int)
		if err0 && err1 && err2 {
			fmt.Println("1003 reply ", arg0, arg1, arg2)
		}
	} else {
		fmt.Println("1003 reply ", err)
	}
}

func run_in_server(x int, str string) {
	fmt.Println("x:", x, "str:", str)
}

func TestServer(t *testing.T) {
	var test_server = TestModule{}
	Start(&test_server)
	time.Sleep(time.Second) //等待goroutine启动
	Call("TestModule", func() {
		run_in_server(777, "run_in_server")
	})

	test_send_msg()
	SyncClose()
}
