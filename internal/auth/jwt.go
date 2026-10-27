package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uint   `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret    []byte
	issuer    string
	ttl       time.Duration
}

func NewJWTService(secret string, ttlMinutes int) *JWTService {
	if ttlMinutes <= 0 {
		ttlMinutes = 60
	}
	return &JWTService{
		secret: []byte(secret),
		issuer: "icecream-inventory",
		ttl:    time.Duration(ttlMinutes) * time.Minute,
	}
}

func (s *JWTService) generateJTI() string {
	// Prefer crypto/rand hex for simplicity
	b := make([]byte, 32)
	if _, err := rand.Read(b); err == nil {
		return hex.EncodeToString(b)
	}
	// Fallback
	return uuid.New().String()
}

func (s *JWTService) Issue(userID uint, role string) (string, string, int64, error) {
	now := time.Now()
	jti := s.generateJTI()
	exp := now.Add(s.ttl)
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ID:        jti,
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", "", 0, err
	}
	return signed, jti, exp.Unix(), nil
}

func (s *JWTService) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	if claims.Issuer != s.issuer {
		return nil, errors.New("invalid issuer")
	}
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}
	if claims.ID == "" {
		return nil, errors.New("missing jti")
	}
	return claims, nil
}
