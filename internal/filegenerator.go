package internal

import (
	"crypto/md5"
	"fmt"
)

type FileSize int

const (
	File100KB FileSize = 100 * 1024
	File1MB   FileSize = 1048576
	File10MB           = 10 * File1MB
	File100MB          = 10 * File10MB
	File1GB   FileSize = 1073741824
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateFileBytes(size FileSize) (buf []byte, hash string) {
	buf = make([]byte, size)
	for i := 0; i < int(size); i++ {
		buf[i] = letters[i%len(letters)]
	}
	return buf, fmt.Sprintf("%x", md5.Sum(buf))
}
