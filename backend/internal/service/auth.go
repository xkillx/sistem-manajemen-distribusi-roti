package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/smdr/backend/internal/config"
	sqlcdb "github.com/smdr/backend/internal/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrMustChangePassword = errors.New("password change required")
	ErrInactiveAccount    = errors.New("account is inactive")
	ErrUsernameExists     = errors.New("username already exists")
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	queries *sqlcdb.Queries
	cfg     *config.Config
}

func NewAuthService(queries *sqlcdb.Queries, cfg *config.Config) *AuthService {
	return &AuthService{queries: queries, cfg: cfg}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, *sqlcdb.User, error) {
	user, err := s.queries.FindUserByUsername(ctx, username)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if !user.Active {
		return "", nil, ErrInactiveAccount
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if user.MustChangePassword {
		return "", nil, ErrMustChangePassword
	}

	token, err := s.generateToken(&user)
	if err != nil {
		return "", nil, fmt.Errorf("generate token: %w", err)
	}

	return token, &user, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.queries.FindUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(oldPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.queries.UpdateUserPassword(ctx, sqlcdb.UpdateUserPasswordParams{
		ID:             userID,
		HashedPassword: string(hashed),
		MustChangePassword: false,
	})
}

func (s *AuthService) ForceChangePassword(ctx context.Context, userID int64, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.queries.UpdateUserPassword(ctx, sqlcdb.UpdateUserPasswordParams{
		ID:             userID,
		HashedPassword: string(hashed),
		MustChangePassword: false,
	})
}

func (s *AuthService) SetTemporaryPassword(ctx context.Context, userID int64) (string, error) {
	tempPass, err := generateTempPassword()
	if err != nil {
		return "", err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(tempPass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	if err := s.queries.UpdateUserPassword(ctx, sqlcdb.UpdateUserPasswordParams{
		ID:             userID,
		HashedPassword: string(hashed),
		MustChangePassword: true,
	}); err != nil {
		return "", err
	}

	return tempPass, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) generateToken(user *sqlcdb.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "smdr",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func generateTempPassword() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
