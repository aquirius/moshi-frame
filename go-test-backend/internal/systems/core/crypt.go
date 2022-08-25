package core

import (
	"crypto/md5"
	"crypto/rand"

	"fmt"
	"io"
)

func (c *Core) generateSessionID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (c *Core) encryptPassword(pw string) string {
	crypt := md5.New()
	io.WriteString(crypt, pw)
	return fmt.Sprintf("%x", crypt.Sum(nil))
}
