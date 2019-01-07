/**
 * Created by angelina on 2018/5/31.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package submail

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/vannnnish/easyweb"
	"github.com/vannnnish/yeego"
	"github.com/vannnnish/yeego/third/yeeSubmail"
	"io/ioutil"
	"net/http"
)

const (
	api                      = "https://api.mysubmail.com/service/verifyphonenumber"
	ProjectIdRegister        = "register"
	ProjectIdForgetPassword  = "forget_password"
	ProjectIdBindPhoneNumber = "bind_phone_number"
	ProjectIdLoginOrRegister = "login_register"
	ProjectIdPreRegister     = "pre_register"

	// 中文
	projectIdRegisterCn        = "ij8cl3" // 注册
	projectIdForgetPasswordCn  = "ij8cl3" // 忘记密码
	projectIdBindPhoneNumberCn = "ij8cl3"  // 绑定手机
	projectIdLoginOrRegisterCn = "ij8cl3" // 登录
	projectIdPreRegisterCn     = "ij8cl3" // 预注册
	// 英文
	projectIdRegisterEn        = "ij8cl3"
	projectIdForgetPasswordEn  = "ij8cl3"
	projectIdBindPhoneNumberEn = "ij8cl3"
	projectIdLoginOrRegisterEn = "ij8cl3"
	projectIdPreRegisterEn     = "ij8cl3"
)

func SendSmsCode(project, areaCode, phoneNum, code string) error {
	projectId := judgeProjectIdByProjectAndAreaCode(project, areaCode)
	var result string
	if areaCode == "86" {
		config := yeeSubmail.Config{
			AppId:    yeego.Config.GetString("submail.AppId"),
			AppKey:   yeego.Config.GetString("submail.AppKey"),
			SignType: "md5",
		}
		mXSend := yeeSubmail.CreateMessageXSend(phoneNum, projectId)
		mXSend.AddVar("code", code)
		result = mXSend.Run(config)
	} else {
		config := yeeSubmail.Config{
			AppId:    yeego.Config.GetString("submail.InternationalAppId"),
			AppKey:   yeego.Config.GetString("submail.InternationalAppKey"),
			SignType: "md5",
		}
		mXSend := yeeSubmail.CreateMessageXSend("+"+areaCode+phoneNum, projectId)
		mXSend.AddVar("code", code)
		result = mXSend.RunInternational(config)
	}
	res := &yeeSubmail.SubmailResponse{}
	if err := json.Unmarshal([]byte(result), res); err != nil {
		easyweb.Logger.Error(result)
		return err
	}
	if res.Status != "success" {
		if res.Code == 115 {
			return errors.New("request-limit")
		}
		return errors.New(res.Msg)
	}
	return nil
}

func judgeProjectIdByProjectAndAreaCode(project, areaCode string) string {
	var projectId string
	if areaCode == "86" {
		switch project {
		case ProjectIdRegister:
			projectId = projectIdRegisterCn
		case ProjectIdForgetPassword:
			projectId = projectIdForgetPasswordCn
		case ProjectIdBindPhoneNumber:
			projectId = projectIdBindPhoneNumberCn
		case ProjectIdLoginOrRegister:
			projectId = projectIdLoginOrRegisterCn
		case ProjectIdPreRegister:
			projectId = projectIdPreRegisterCn
		}
	} else {
		switch project {
		case ProjectIdRegister:
			projectId = projectIdRegisterEn
		case ProjectIdForgetPassword:
			projectId = projectIdForgetPasswordEn
		case ProjectIdBindPhoneNumber:
			projectId = projectIdBindPhoneNumberEn
		case ProjectIdLoginOrRegister:
			projectId = projectIdLoginOrRegisterEn
		case ProjectIdPreRegister:
			projectId = projectIdPreRegisterEn
		}
	}
	return projectId
}

type Result struct {
	Status        string `json:"status"`
	Code          int    `json:"code;omitempty"`
	Msg           string `json:"msg;omitempty"`
	To            string `json:"to"`
	CountryCode   int    `json:"country_code"`
	CountryRegion string `json:"country_region"`
	CountryName   string `json:"country_name"`
	Price         string `json:"price"`
}

// 验证国际手机号码格式
func VerifyInternationalPhoneNumber(areaCode, phoneNum string) (Result, error) {
	to := "+" + areaCode + phoneNum
	m := make(map[string]string)
	m["to"] = to
	data, err := httpPost(api, m)
	if err != nil {
		return Result{}, err
	}
	var res Result
	err = json.Unmarshal(data, &res)
	return res, err
}

func httpPost(queryUrl string, postData map[string]string) ([]byte, error) {
	data, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer([]byte(data))
	retStr, err := http.Post(queryUrl, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(retStr.Body)
	retStr.Body.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}
