package utils

import (
	"crypto/rand"
	"fmt"
)

type UID string // represent of uuid identificator

//GenerateUID generate unique identificator
func GenerateUID() (UID, error) {
	var uid UID

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return uid, err
	}

	uid = UID(fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))

	return uid, nil
}
