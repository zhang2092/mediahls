package password

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// ******************** scrypt ********************

// ScryptHashPassword scrypt 加密
// password 原始密码
func ScryptHashPassword(password string) (string, error) {
	// example for making salt - https://play.golang.org/p/_Aw6WeWC42I
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// using recommended cost parameters from - https://godoc.org/golang.org/x/crypto/scrypt
	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// return hex-encoded string with salt appended to password
	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

// ScryptComparePassword 判断密码是否正确
// storedPassword 加密密码
// suppliedPassword 原始密码
func ScryptComparePassword(storedPassword string, suppliedPassword string) error {
	pwsalt := strings.Split(storedPassword, ".")

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwsalt[1])

	if err != nil {
		return fmt.Errorf("unable to verify user password")
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return err
	}

	if hex.EncodeToString(shash) == pwsalt[0] {
		return nil
	}

	return fmt.Errorf("password error")
}

// ******************** bcrypt ********************

// BcryptHashPassword bcrypt 加密
// password 原始密码
func BcryptHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// BcryptComparePassword 判断密码是否正确
// hashedPassword 加密密码
// password 原始密码
func BcryptComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
