package utils

import (
	"crypto/md5"
	"fmt"
)

func Md5String(str string) string {
	md5Str := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", md5Str)
}
