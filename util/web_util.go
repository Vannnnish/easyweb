/*
@Time : 2019-01-06 23:15 
@Author : vannnnish
@File : web_util
*/

// web常用到的一些功能函数
package util

import (
	"errors"
	"github.com/vannnnish/easyweb/submail"
	"github.com/vannnnish/yeego/yeeStrconv"
	"github.com/vannnnish/yeego/yeeStrings"
	"regexp"
	"strings"
)

// 去除区号前面的+号
func FormatAreaCode(areaCode string) string {
	if strings.HasPrefix(areaCode, "+") {
		areaCode = strings.Replace(areaCode, "+", "", -1)
	}
	return areaCode
}

var areaCodeArr = []string{"93", "355", "213", "684", "376", "244", "1264", "1268", "54", "374", "297", "247", "61", "43", "994", "1242", "973", "880", "1246", "375", "32", "501", "229", "1441", "975", "591", "387", "267", "55", "1284", "673", "359", "226", "257", "855", "237", "1", "238", "1345", "236", "235", "56", "86", "57", "269", "242", "243", "682", "506", "225", "385", "53", "357", "420", "243", "45", "253", "1767", "1809", "670", "593", "20", "503", "240", "291", "372", "251", "500", "298", "679", "358", "33", "594", "689", "809", "245", "241", "220", "995", "49", "233", "350", "30", "299", "1473", "590", "1671", "502", "44", "224", "592", "509", "504", "852", "36", "354", "91", "62", "98", "964", "353", "44", "972", "39", "1876", "81", "44", "962", "7", "254", "965", "996", "856", "371", "961", "266", "231", "218", "423", "370", "352", "853", "389", "261", "265", "60", "960", "223", "356", "596", "222", "230", "52", "691", "373", "377", "976", "382", "1664", "212", "258", "264", "977", "31", "599", "687", "64", "505", "227", "234", "61", "1", "47", "968", "92", "680", "970", "507", "675", "595", "51", "63", "48", "351", "1787", "974", "95", "262", "40", "7", "250", "508", "1809", "1758", "784", "378", "239", "966", "221", "381", "248", "232", "65", "421", "386", "677", "252", "27", "82", "211", "34", "94", "1758", "1784", "1869", "249", "597", "268", "46", "41", "963", "886", "992", "255", "66", "228", "676", "1868", "216", "90", "993", "1649", "688", "1340", "256", "380", "971", "44", "1", "598", "998", "678", "58", "84", "685", "967", "260", "263",}

// 验证区号
func ValidAreaCode(areaCode string) bool {
	areaCode = FormatAreaCode(areaCode)
	return yeeStrings.IsInSlice(areaCodeArr, areaCode)
}

var lanArr = []string{"zh-cn", "en-us"}

// 验证语言环境
func ValidLanguage(lan string) bool {
	return yeeStrings.IsInSlice(lanArr, lan)
}

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")
var mobilePattern = regexp.MustCompile("^((\\+86)|(86))?1([38]\\d|4[579]|5[0-35-9]|7[0135-8])\\d{8}$")

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

// 验证国内/国际手机号码格式
func ValidPhoneNumber(areaCode, phoneNumber string) (bool, error) {
	if areaCode == "86" {
		return mobilePattern.Match([]byte(phoneNumber)), nil
	}
	res, err := submail.VerifyInternationalPhoneNumber(areaCode, phoneNumber)
	if err != nil {
		// 为什么返回true呢，因为这个地方说明是请求失败了，那么不知道是否号码正确，所以认定正确吧~~~
		return true, err
	}
	if res.Status == "success" {
		return true, nil
	}
	return false, errors.New(yeeStrconv.FormatInt(res.Code) + " : " + res.Msg)
}

// 验证邮箱格式
func ValidEmail(email string) bool {
	return emailPattern.Match([]byte(email))
}
