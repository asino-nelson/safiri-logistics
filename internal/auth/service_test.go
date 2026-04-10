package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asino-nelson/safiri-logistics/internal/user"
)

func TestServiceRegister(t *testing.T) {
	t.Parallel()

	repo := newStubUserRepository()
	service := NewService(repo, BcryptHasher{}, NewJWTManager("secret", time.Hour))

	response, err := service.Register(context.Background(), RegisterRequest{
		Name:     "Shipper One",
		Email:    "shipper@example.com",
		Password: "password123",
		Role:     user.RoleCustomer,
	})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}

	if response.User.Email != "shipper@example.com" {
		t.Fatalf("expected normalized email, got %s", response.User.Email)
	}

	stored, err := repo.GetByEmail(context.Background(), "shipper@example.com")
	if err != nil {
		t.Fatalf("expected stored user: %v", err)
	}

	if stored.PasswordHash == "password123" {
		t.Fatal("expected hashed password to differ from plain text")
	}
}

func TestServiceLoginRejectsInvalidPassword(t *testing.T) {
	t.Parallel()

	hasher := BcryptHasher{}
	hash, err := hasher.Hash("correct-password")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := newStubUserRepository()
	repo.users["driver@example.com"] = &user.User{
		ID:           "user-1",
		Name:         "Driver",
		Email:        "driver@example.com",
		PasswordHash: hash,
		Role:         user.RoleDriver,
	}

	service := NewService(repo, hasher, NewJWTManager("secret", time.Hour))

	_, err = service.Login(context.Background(), LoginRequest{
		Email:    "driver@example.com",
		Password: "bad-password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

type stubUserRepository struct {
	users map[string]*user.User
}

func newStubUserRepository() *stubUserRepository {
	return &stubUserRepository{
		users: make(map[string]*user.User),
	}
}

func (r *stubUserRepository) Create(_ context.Context, account *user.User) error {
	if _, exists := r.users[account.Email]; exists {
		return user.ErrEmailAlreadyUsed
	}

	clone := *account
	r.users[clone.Email] = &clone
	return nil
}

func (r *stubUserRepository) GetByEmail(_ context.Context, email string) (*user.User, error) {
	account, exists := r.users[email]
	if !exists {
		return nil, user.ErrUserNotFound
	}

	clone := *account
	return &clone, nil
}

func (r *stubUserRepository) GetByID(_ context.Context, userID string) (*user.User, error) {
	for _, account := range r.users {
		if account.ID == userID {
			clone := *account
			return &clone, nil
		}
	}

	return nil, user.ErrUserNotFound
}
