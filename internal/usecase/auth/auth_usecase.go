package auth_usercase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/DannyAss/users/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID  uint   `json:"uid"`
	Name    string `json:"name"`
	Email   string `json:"email,omitempty"`
	RoleIds uint   `json:"role_ids,omitempty"`
	jwt.RegisteredClaims
}

type AuthUsecase interface {
}

type authUsecase struct {
	cfg *config.ConfigEnv
}

func NewAuthUsecase(cfg *config.ConfigEnv) AuthUsecase {
	return &authUsecase{cfg: cfg}
}

func GenerateAccessToken(secret []byte, userID uint, email, name string, role_ids uint, ttl time.Duration) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(ttl)
	claims := JWTClaims{
		UserID:  userID,
		Name:    name,
		Email:   email,
		RoleIds: role_ids,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)

	return signed, exp, err
}

func VerifyAccessToken(secret []byte, tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			fmt.Println("cek 1")
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		fmt.Println("cek 2")
		fmt.Println(err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	fmt.Println("cek 3")
	return nil, errors.New("invalid token")
}

func GenerateRefreshToken() (token string, raw []byte, err error) {
	raw = make([]byte, 32) // 256-bit
	_, err = rand.Read(raw)
	if err != nil {
		return "", nil, err
	}

	enc := base64.RawURLEncoding.EncodeToString(raw)
	return enc, raw, nil
}
