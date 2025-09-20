package token

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义JWT claims结构体
type Claims struct {
	jwt.RegisteredClaims
	UserID uint64 `json:"user_id"`
	UType  int    `json:"type"` // 区分用户类型
}

func (c Claims) IsAdmin() bool {
	return c.UType == 1
}

func (c Claims) String() string {
	jsonBytes, _ := json.Marshal(c)
	return string(jsonBytes)
}

// NewClaims 创建一个新的Claims实例
func NewAdminClaims(userID uint64, expiration time.Duration) Claims {
	now := time.Now()
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID: userID,
		UType:  1,
	}
}

// NewClaims 创建一个新的Claims实例
func NewClaims(userID uint64, expiration time.Duration) Claims {
	now := time.Now()
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID: userID,
	}
}

// TokenPair 包含访问令牌和刷新令牌的结构体
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
