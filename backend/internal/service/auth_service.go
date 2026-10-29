package service

import (
	"errors"
	"strings"
	"time"

	"github.com/MrVirul/icecream-inventory/internal/auth"
	"github.com/MrVirul/icecream-inventory/internal/httpx"
	"github.com/MrVirul/icecream-inventory/internal/models"
	"github.com/MrVirul/icecream-inventory/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (*LoginResult, error)
	Logout(jti string, expUnix int64) error
}

type LoginResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
	User        *models.User
}

type authService struct {
	users      repository.UserRepository
	revoked    repository.RevokedTokenRepository
	jwt        *auth.JWTService
	bcryptCost int
}

func NewAuthService(users repository.UserRepository, revoked repository.RevokedTokenRepository, jwtSvc *auth.JWTService, bcryptCost int) AuthService {
	if bcryptCost < 4 {
		bcryptCost = 12
	}
	return &authService{users: users, revoked: revoked, jwt: jwtSvc, bcryptCost: bcryptCost}
}

func (s *authService) Login(email, password string) (*LoginResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	user, err := s.users.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Equalize timing
			_ = bcrypt.CompareHashAndPassword([]byte("$2a$12$C6UzMDM.H6dfI/f/IKcEe.OXGXMfPYM77ZJ4KAVTOhW6/fePq6v32"), []byte(password))
			return nil, httpx.Unauthorized("invalid email or password")
		}
		return nil, httpx.Unauthorized("invalid email or password")
	}
	if err := auth.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, httpx.Unauthorized("invalid email or password")
	}
	token, jti, expUnix, err := s.jwt.Issue(user.ID, user.Role)
	if err != nil {
		return nil, httpx.InternalError()
	}
	// Optionally store jti mapping not required if we just check revoked table on auth; but we have jti in token
	_ = jti
	expiresIn := expUnix - time.Now().Unix()
	if expiresIn < 0 {
		expiresIn = 0
	}
	return &LoginResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User:        user,
	}, nil
}

func (s *authService) Logout(jti string, expUnix int64) error {
	if jti == "" {
		return nil
	}
	exp := time.Unix(expUnix, 0)
	if exp.IsZero() {
		exp = time.Now().Add(time.Hour)
	}
	return s.revoked.Create(jti, exp)
}
