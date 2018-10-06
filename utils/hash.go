package utils

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

// Hash sha3 512 哈希
func Hash(s string) string {
	return fmt.Sprintf("%x", sha3.Sum512([]byte(s)))
}
