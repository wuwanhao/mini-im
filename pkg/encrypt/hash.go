package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// MD5 加密
func MD5(str []byte) string {
	hash := md5.New()
	hash.Write(str)
	return hex.EncodeToString(hash.Sum(nil))
}

// 使用 bcrypt 生成密码哈希
func GenPasswordHash(password []byte) ([]byte, error) {
	// 生成哈希
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// 验证密码和哈希是否匹配
func ValidatePasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
