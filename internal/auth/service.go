package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidRole        = errors.New("invalid role")
)

type UserRepository interface {
	Create(ctx context.Context, user *user.User) error
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenManager interface {
	Generate(userID, role string) (string, error)
	Parse(token string) (*TokenClaims, error)
}

type Service struct {
	users   UserRepository
	hasher  PasswordHasher
	tokens  TokenManager
	nowFunc func() time.Time
}

func NewService(users UserRepository, hasher PasswordHasher, tokens TokenManager) *Service {
	return &Service{
		users:   users,
		hasher:  hasher,
		tokens:  tokens,
		nowFunc: time.Now,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterRequest) (*AuthResponse, error) {
	role := normalizeRole(input.Role)
	if role != user.RoleCustomer && role != user.RoleDriver {
		return nil, ErrInvalidRole
	}

	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	account := &user.User{
		ID:           uuid.NewString(),
		Name:         strings.TrimSpace(input.Name),
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    s.nowFunc().UTC(),
	}

	if err := s.users.Create(ctx, account); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(account)
}

func (s *Service) Login(ctx context.Context, input LoginRequest) (*AuthResponse, error) {
	account, err := s.users.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := s.hasher.Compare(account.PasswordHash, input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.buildAuthResponse(account)
}

func (s *Service) buildAuthResponse(account *user.User) (*AuthResponse, error) {
	token, err := s.tokens.Generate(account.ID, account.Role)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User: UserPayload{
			ID:    account.ID,
			Name:  account.Name,
			Email: account.Email,
			Role:  account.Role,
		},
	}, nil
}

type BcryptHasher struct{}

func (BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type JWTManager struct {
	secret   []byte
	lifetime time.Duration
}

func NewJWTManager(secret string, lifetime time.Duration) *JWTManager {
	return &JWTManager{
		secret:   []byte(secret),
		lifetime: lifetime,
	}
}

func (m *JWTManager) Generate(userID, role string) (string, error) {
	now := time.Now().UTC()
	claims := TokenClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.lifetime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Parse(token string) (*TokenClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(_ *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok || !parsedToken.Valid {
		return nil, ErrInvalidCredentials
	}

	return claims, nil
}

func normalizeRole(role string) string {
	return strings.ToLower(strings.TrimSpace(role))
}
