package utils

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"
)

func Md5Crypt(str string, salt ...interface{}) (CryptStr string) {
	if l := len(salt); l > 0 {
		slice := make([]string, l+1)
		str = fmt.Sprintf(str+strings.Join(slice, "%v"), salt...)
	}
	log.Println("Md5Crypt str:", str)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
