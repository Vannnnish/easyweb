/*
@Time : 2019-01-07 14:41 
@Author : vannnnish
@File : json_tool
*/

package json_tool

import "github.com/json-iterator/go"


var json = jsoniter.Config{
	EscapeHTML:                    false,
	SortMapKeys:                   false,
	ValidateJsonRawMessage:        false,
	ObjectFieldMustBeSimpleString: true,
	MarshalFloatWith6Digits:       false,
}.Froze()

func Json() jsoniter.API {
	return json
}

func MarshalToString(v interface{}) string {
	str, _ := json.MarshalToString(v)
	return str
}

func Marshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}