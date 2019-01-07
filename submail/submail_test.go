/**
 * Created by angelina on 2018/6/13.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package submail

import (
	"github.com/vannnnish/easyweb"
	"github.com/vannnnish/easyweb/log4go"
	"github.com/vannnnish/yeego"
	"testing"
)

func TestVerifyInternationalPhoneNumber(t *testing.T) {
	p1 := "15008477531"
	res, err := VerifyInternationalPhoneNumber("1", p1)
	if err != nil {
		panic(err.Error())
	}
	yeego.Print(res)
}

func TestSendSmsCode(t *testing.T) {
	yeego.MustInitConfig("../conf", "conf")
	easyweb.Logger = log4go.NewDefaultLogger(log4go.ERROR)
	err := SendSmsCode(ProjectIdRegister, "86", "17318070950", "123456")
	if err != nil {
		panic(err.Error())
	}
}
