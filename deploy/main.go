package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05Z07:00"
)

// sha256加密
func Sha256(data []byte) string {
	_sha1 := sha256.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}
func main() {
	now := time.Now().Format(timeLayout)
	fmt.Println(now)
}
