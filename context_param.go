/**
 * Created by angelina on 2017/8/28.
 */

package easyweb

import (
	"errors"
	"github.com/buger/jsonparser"
	"github.com/spf13/cast"
	"github.com/vannnnish/yeego/yeecrypto/aes"
	"github.com/vannnnish/yeego/yeecrypto/rsa"
	"github.com/vannnnish/yeego/yeestrconv"
	"strconv"
)

// 链式调用获取参数
func (c *Context) ClearParam() {
	c.nowParam = *new(Param)
}

func (c *Context) GetParam(key string) *Context {
	c.ClearParam()
	c.nowParam.Key = key
	c.nowParam.Value = c.Query(key)
	return c
}

func (c *Context) PostParam(key string) *Context {
	c.ClearParam()
	c.nowParam.Key = key
	c.nowParam.Value = c.PostForm(key)
	return c
}

func (c *Context) Param(key string) *Context {
	c.ClearParam()
	c.nowParam.Key = key
	if c.Query(key) != "" {
		c.nowParam.Value = c.Query(key)
	} else {
		c.nowParam.Value = c.PostForm(key)
	}
	return c
}

func (c *Context) SetDefault(val string) *Context {
	if len(c.nowParam.Value) == 0 {
		c.nowParam.Value = val
	}
	return c
}

func (c *Context) SetDefaultInt(i int) *Context {
	if len(c.nowParam.Value) == 0 {
		c.nowParam.Value = yeestrconv.FormatInt(i)
	}
	return c
}

func (c *Context) GetString() string {
	return c.nowParam.Value
}

func (c *Context) MustGetString() string {
	if len(c.nowParam.Value) == 0 {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, "参数不能为空,参数名称为:"+c.nowParam.Key)
		}
	}
	return c.nowParam.Value
}

func (c *Context) MustGetStringWithError(str string) string {
	if len(c.nowParam.Value) == 0 {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, str)
		}
	}
	return c.nowParam.Value
}

func (c *Context) GetInt() int {
	if len(c.nowParam.Value) == 0 {
		return 0
	}
	i, err := strconv.Atoi(c.nowParam.Value)
	if err != nil {
		return 0
	}
	return i
}

func (c *Context) MustGetInt() int {
	if len(c.nowParam.Value) == 0 {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, "参数不能为空,参数名称为:"+c.nowParam.Key)
		}
		return 0
	}
	i, err := strconv.Atoi(c.nowParam.Value)
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, "参数类型错误,参数名称为:"+c.nowParam.Key)
		}
		return 0
	}
	return i
}

func (c *Context) MustGetIntWithError(str string) int {
	if len(c.nowParam.Value) == 0 {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, str)
		}
		return 0
	}
	i, err := strconv.Atoi(c.nowParam.Value)
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, str)
		}
		return 0
	}
	return i
}

func (c *Context) GetBool() bool {
	if len(c.nowParam.Value) == 0 {
		return false
	}
	value, err := strconv.ParseBool(c.nowParam.Value)
	if err != nil {
		return false
	}
	return value
}

func (c *Context) MustGetBool() bool {
	if len(c.nowParam.Value) == 0 {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, "参数不能为空,参数名称为:"+c.nowParam.Key)
		}
		return false
	}
	value, err := strconv.ParseBool(c.nowParam.Value)
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, "参数类型错误,参数名称为:"+c.nowParam.Key)
		}
		return false
	}
	return value
}

func (c *Context) MustGetBoolWithError(str string) bool {
	if len(c.nowParam.Value) == 0 {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, str)
		}
		return false
	}
	value, err := strconv.ParseBool(c.nowParam.Value)
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, str)
		}
		return false
	}
	return value
}

func (c *Context) GetFloat() float64 {
	if len(c.nowParam.Value) == 0 {
		return 0
	}
	f, err := strconv.ParseFloat(c.nowParam.Value, 64)
	if err != nil {
		return 0
	}
	return f
}

func (c *Context) MustGetFloat() float64 {
	if len(c.nowParam.Value) == 0 {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, "参数不能为空,参数名称为:"+c.nowParam.Key)
		}
		return 0
	}
	f, err := strconv.ParseFloat(c.nowParam.Value, 64)
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, "参数类型错误,参数名称为:"+c.nowParam.Key)
		}
		return 0
	}
	return f
}

func (c *Context) MustGetFloatWithError(str string) float64 {
	if len(c.nowParam.Value) == 0 {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, str)
		}
		return 0
	}
	f, err := strconv.ParseFloat(c.nowParam.Value, 64)
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(c.nowParam.Key, str)
		}
		return 0
	}
	return f
}

func (c *Context) GetError() error {
	if c.validation.HasErrors() {
		for _, err := range c.validation.Errors {
			return errors.New(err.Message)
		}
	}
	return nil
}

/**************************RSA  &  AES*******************************/

const (
	// 服务端解密客户端数据使用的私钥
	ServerPrivateKeyPath string = "cert/server_private.pem"
	// 服务端加密数据使用的公钥
	ClientPublicKeyPath string = "cert/client_public.pem"
)

func (c *Context) RsaGetString(key string) string {
	data := c.Param("data").GetString()
	jsonBytes, err := rsa.RsaDecrypt(ServerPrivateKeyPath, []byte(data))
	if err != nil {
		return ""
	}
	str, err := jsonparser.GetString(jsonBytes, key)
	if err != nil {
		return ""
	}
	return str
}

func (c *Context) RsaGetInt(key string, defaultInt int) int {
	data := c.Param("data").GetString()
	jsonBytes, err := rsa.RsaDecrypt(ServerPrivateKeyPath, []byte(data))
	if err != nil {
		return defaultInt
	}
	i, err := jsonparser.GetInt(jsonBytes, key)
	if err != nil {
		return defaultInt
	}
	return int(i)
}

func (c *Context) RsaMustGetString(key string, errStr string) string {
	data := c.Param("data").GetString()
	jsonBytes, err := rsa.RsaDecrypt(ServerPrivateKeyPath, []byte(data))
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(key, errStr)
		}
		return ""
	}
	str, err := jsonparser.GetString(jsonBytes, key)
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(key, errStr)
		}
		return ""
	}
	return str
}

func (c *Context) RsaMustGetInt(key string, errStr string) int {
	data := c.Param("data").GetString()
	jsonBytes, err := rsa.RsaDecrypt(ServerPrivateKeyPath, []byte(data))
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(key, errStr)
		}
		return 0
	}
	i, err := jsonparser.GetInt(jsonBytes, key)
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(key, errStr)
		}
		return 0
	}
	return int(i)
}

func (c *Context) AesGetInt(param string, defaultInt int) int {
	data := c.Param("data").MustGetString()
	if data == "" {
		// 未登录
		return c.Param(param).SetDefaultInt(defaultInt).GetInt()
	}
	decryptData, err := aes.AesDecrypt([]byte(c.StoreGetString("key")), []byte(data))
	if err != nil {
		return defaultInt
	}
	i, err := jsonparser.GetInt(decryptData, param)
	if err != nil {
		iStr, err := jsonparser.GetString(decryptData, param)
		if err != nil {
			return defaultInt
		}
		i64, err := strconv.Atoi(iStr)
		if err != nil {
			return defaultInt
		}
		return int(i64)
	}
	return int(i)
}

func (c *Context) AesMustGetIntWithError(param string, errStr string) int {
	data := c.Param("data").MustGetString()
	if data == "" {
		// 未登录
		return c.Param(param).MustGetIntWithError(errStr)
	}
	decryptData, err := aes.AesDecrypt([]byte(c.StoreGetString("key")), []byte(data))
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(param, errStr)
		}
		return 0
	}
	i, err := jsonparser.GetInt(decryptData, param)
	if err != nil {
		iStr, err := jsonparser.GetString(decryptData, param)
		if err != nil {
			if !c.validation.HasErrors() {
				c.validation.SetError(param, errStr)
			}
			return 0
		}
		i64, err := strconv.Atoi(iStr)
		if err != nil {
			if !c.validation.HasErrors() {
				c.validation.SetError(param, errStr)
			}
			return 0
		}
		return int(i64)
	}
	return int(i)
}

func (c *Context) AesGetString(param string, defaultStr ...string) string {
	data := c.Param("data").MustGetString()
	if data == "" {
		return c.Param(param).GetString()
	}
	decryptData, err := aes.AesDecrypt([]byte(c.StoreGetString("key")), []byte(data))
	if err != nil {
		return ""
	}
	str, err := jsonparser.GetString(decryptData, param)
	if err != nil || len(str) == 0 {
		if len(defaultStr) > 0 {
			return defaultStr[0]
		}
		return ""
	}
	return str
}

func (c *Context) AesMustGetStringWithError(param, errStr string) string {
	data := c.Param("data").MustGetString()
	if data == "" {
		return c.Param(param).MustGetStringWithError(errStr)
	}
	decryptData, err := aes.AesDecrypt([]byte(c.StoreGetString("key")), []byte(data))
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(param, errStr)
		}
		return ""
	}
	str, err := jsonparser.GetString(decryptData, param)
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(param, errStr)
		}
		return ""
	}
	return str
}

func (c *Context) AesGetFloat(param string) float64 {
	data := c.Param("data").MustGetString()
	if data == "" {
		return c.Param(param).GetFloat()
	}
	decryptData, err := aes.AesDecrypt([]byte(c.StoreGetString("key")), []byte(data))
	if err != nil {
		return 0
	}
	f, err := jsonparser.GetFloat(decryptData, param)
	if err != nil {
		fStr, err := jsonparser.GetString(decryptData, param)
		if err != nil {
			return 0
		}
		f, err := yeestrconv.ParseFloat64(fStr)
		if err != nil {
			return 0
		}
		return f
	}
	return f
}

func (c *Context) AesMustGetFloatWithError(param, errStr string) float64 {
	data := c.Param("data").MustGetString()
	if data == "" {
		return c.Param(param).MustGetFloatWithError(errStr)
	}
	decryptData, err := aes.AesDecrypt([]byte(c.StoreGetString("key")), []byte(data))
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(param, errStr)
		}
		return 0
	}
	f, err := jsonparser.GetFloat(decryptData, param)
	if err != nil {
		fStr, err := jsonparser.GetString(decryptData, param)
		if err != nil {
			if !c.validation.HasErrors() {
				c.validation.SetError(param, errStr)
			}
			return 0
		}
		f, err := yeestrconv.ParseFloat64(fStr)
		if err != nil {
			if !c.validation.HasErrors() {
				c.validation.SetError(param, errStr)
			}
			return 0
		}
		return f
	}
	return f
}

func (c *Context) AesGetBool(param string, defaultB bool) bool {
	data := c.Param("data").MustGetString()
	if data == "" {
		return c.Param(param).GetBool()
	}
	decryptData, err := aes.AesDecrypt([]byte(c.StoreGetString("key")), []byte(data))
	if err != nil {
		return defaultB
	}
	//b, err := jsonparser.GetBoolean(decryptData, param)
	bBytes, _, _, err1 := jsonparser.Get(decryptData, param)
	b, err2 := strconv.ParseBool(cast.ToString(bBytes))
	if err1 != nil || err2 != nil {
		return defaultB
	}
	return b
}

func (c *Context) AesMustGetBool(param, errStr string) bool {
	data := c.Param("data").MustGetString()
	if data == "" {
		return c.Param(param).MustGetBoolWithError(errStr)
	}
	decryptData, err := aes.AesDecrypt([]byte(c.StoreGetString("key")), []byte(data))
	if err != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(param, errStr)
		}
		return false
	}
	bBytes, _, _, err1 := jsonparser.Get(decryptData, param)
	b, err2 := strconv.ParseBool(cast.ToString(bBytes))
	if err1 != nil || err2 != nil {
		if !c.validation.HasErrors() {
			c.validation.SetError(param, errStr)
		}
		return false
	}
	return b
}
