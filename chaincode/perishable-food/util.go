package main

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
)

// sha256加密
func Sha256(data []byte) string {
	_sha1 := sha256.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

// 判断slice是否存在某值 [反射]
func IsExistItem(value interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() { //  判断类型
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}
