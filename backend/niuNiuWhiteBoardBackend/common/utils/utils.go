package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
)

// sha1加密
func Sha1En(data string) string {
	t := sha1.New() ///产生一个散列值得方式
	_, _ = io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
