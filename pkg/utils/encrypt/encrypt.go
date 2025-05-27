package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
)

const Salt = "1VEryjia20n5DanDE271mnnIm25aSa1tq0Uj10"

// EncryptPassword 使用 SHA-256 和多次迭代加密密码
func EncryptPassword(password string) string {
	hash := sha256.New()
	saltedPassword := password + Salt
	for i := 0; i < 1000; i++ { // 进行 1000 次迭代
		hash.Write([]byte(saltedPassword))
		saltedPassword = hex.EncodeToString(hash.Sum(nil))
		hash.Reset()
	}
	return saltedPassword
}
