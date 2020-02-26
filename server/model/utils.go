package model

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
)

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")

func NewRandomString(length int) string {
	var b bytes.Buffer
	str := make([]byte, length+8)
	rand.Read(str)
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(str)
	encoder.Close()
	b.Truncate(length) // removes the '==' padding
	return b.String()
}
