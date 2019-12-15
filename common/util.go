/*
* @Author:	Payton
* @Date:	util
* @DESC: 	DESC
 */

package util

import (
	"gen_log"
)

type CallbackFunc func()

func RunFuncWithPainc(Fun CallbackFunc) {
	defer func() {
		err := recover()
		if err != nil {
			gen_log.Error(err)
		}
	}()

	Fun()
}
