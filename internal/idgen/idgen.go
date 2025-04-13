package idgen

import (
	"crypto/rand"
)

func NewUserID() string {
	return "U_" + newIdentifier(32)
}

func NewTaskID() string {
	return "T_" + newIdentifier(32)
}

func newIdentifier(count int) string {
	const prettyChars = "ABCDEFGHJKLMNPQRSTUVWXYZ0123456789"
	const charsetLen = byte(len(prettyChars))
	idBytes := make([]byte, count)
	if _, err := rand.Read(idBytes); err != nil {
		panic(err)
	}

	for i := range count {
		idBytes[i] = prettyChars[idBytes[i]%charsetLen]
	}

	return string(idBytes)
}
