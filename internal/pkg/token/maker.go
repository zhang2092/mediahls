package token

import (
	"time"
)

// Maker 管理token的接口定义
type Maker interface {
	// CreateToken 根据用户名和时间创建一个新的token
	CreateToken(id string, username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken 校验token是否正确
	VerifyToken(token string) (*Payload, error)
}
