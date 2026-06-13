package services

import (
	"context"
	"testing"
	"time"

	"github.com/Lynccs/payment-service/internal/app/dto"
	"github.com/Lynccs/payment-service/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// --- mock repos ---

type mockUserRepo struct {
	users  map[string]models.User
	nextID int
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]models.User), nextID: 1}
}

func (m *mockUserRepo) Create(_ context.Context, user models.User) (int, error) {
	id := m.nextID
	m.nextID++
	user.ID = id
	user.CreatedAt = time.Now()
	m.users[user.Email] = user
	return id, nil
}

func (m *mockUserRepo) GetByID(_ context.Context, id int) (models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return models.User{}, ErrUserNotFound
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (models.User, error) {
	u, ok := m.users[email]
	if !ok {
		return models.User{}, ErrUserNotFound
	}
	return u, nil
}

type mockWalletRepo struct{}

func (m *mockWalletRepo) Create(_ context.Context, userID int) (models.Wallet, error) {
	return models.Wallet{ID: userID, UserID: userID}, nil
}
func (m *mockWalletRepo) GetByUserID(_ context.Context, _ int) (models.Wallet, error) {
	return models.Wallet{}, nil
}
func (m *mockWalletRepo) GetBalance(_ context.Context, _ int) (float64, error) { return 0, nil }
func (m *mockWalletRepo) SetInitialBalance(_ context.Context, _ int, _ float64) error {
	return nil
}

// --- helpers ---

func newSvc() (*userService, *mockUserRepo) {
	ur := newMockUserRepo()
	svc := &userService{userRepo: ur, walletRepo: &mockWalletRepo{}}
	return svc, ur
}

func registerUser(t *testing.T, svc *userService, email, password string) dto.UserResponse {
	t.Helper()
	resp, err := svc.Register(context.Background(), dto.RegisterRequest{Email: email, Password: password})
	require.NoError(t, err)
	return resp
}

// --- validatePassword ---

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"too short", "Ab1!", ErrPasswordTooShort},
		{"no uppercase", "abcdef1!", ErrPasswordNoUpper},
		{"no digit", "Abcdefg!", ErrPasswordNoDigit},
		{"no symbol", "Abcdefg1", ErrPasswordNoSymbol},
		{"valid", "Abcdef1!", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePassword(tc.password)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

// --- Register ---

func TestRegister_Success(t *testing.T) {
	svc, _ := newSvc()

	resp, err := svc.Register(context.Background(), dto.RegisterRequest{
		Email:    "user@example.com",
		Password: "Secret1!",
	})

	require.NoError(t, err)
	assert.Equal(t, "user@example.com", resp.Email)
	assert.NotZero(t, resp.ID)
	assert.NotZero(t, resp.CreatedAt)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	svc, _ := newSvc()
	registerUser(t, svc, "dup@example.com", "Secret1!")

	_, err := svc.Register(context.Background(), dto.RegisterRequest{
		Email:    "dup@example.com",
		Password: "Secret1!",
	})

	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
}

func TestRegister_PasswordValidation(t *testing.T) {
	cases := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"too short", "Ab1!", ErrPasswordTooShort},
		{"no uppercase", "abcdef1!", ErrPasswordNoUpper},
		{"no digit", "Abcdefg!", ErrPasswordNoDigit},
		{"no symbol", "Abcdefg1", ErrPasswordNoSymbol},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, _ := newSvc()
			_, err := svc.Register(context.Background(), dto.RegisterRequest{
				Email:    "user@example.com",
				Password: tc.password,
			})
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestRegister_PasswordNotStoredInResponse(t *testing.T) {
	svc, ur := newSvc()
	registerUser(t, svc, "user@example.com", "Secret1!")

	stored := ur.users["user@example.com"]
	assert.NotEqual(t, "Secret1!", stored.PasswordHash, "пароль має зберігатись як bcrypt-хеш")
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte("Secret1!")))
}

// --- Login ---

func TestLogin_Success(t *testing.T) {
	svc, _ := newSvc()
	registerUser(t, svc, "user@example.com", "Secret1!")

	resp, err := svc.Login(context.Background(), dto.LoginRequest{
		Email:    "user@example.com",
		Password: "Secret1!",
	})

	require.NoError(t, err)
	assert.Equal(t, "user@example.com", resp.Email)
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, _ := newSvc()

	_, err := svc.Login(context.Background(), dto.LoginRequest{
		Email:    "nobody@example.com",
		Password: "Secret1!",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _ := newSvc()
	registerUser(t, svc, "user@example.com", "Secret1!")

	_, err := svc.Login(context.Background(), dto.LoginRequest{
		Email:    "user@example.com",
		Password: "Wrong1!X",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_WrongAndCorrectPasswordReturnSameError(t *testing.T) {
	svc, _ := newSvc()

	_, errNotFound := svc.Login(context.Background(), dto.LoginRequest{
		Email: "nobody@example.com", Password: "Secret1!",
	})

	registerUser(t, svc, "user@example.com", "Secret1!")
	_, errWrongPass := svc.Login(context.Background(), dto.LoginRequest{
		Email: "user@example.com", Password: "Wrong1!X",
	})

	// обидві помилки мають бути ErrInvalidCredentials — захист від перебору
	assert.ErrorIs(t, errNotFound, ErrInvalidCredentials)
	assert.ErrorIs(t, errWrongPass, ErrInvalidCredentials)
}

// --- GetProfile ---

func TestGetProfile_Success(t *testing.T) {
	svc, _ := newSvc()
	created := registerUser(t, svc, "user@example.com", "Secret1!")

	resp, err := svc.GetProfile(context.Background(), created.ID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, resp.ID)
	assert.Equal(t, "user@example.com", resp.Email)
}

func TestGetProfile_NotFound(t *testing.T) {
	svc, _ := newSvc()

	_, err := svc.GetProfile(context.Background(), 9999)

	assert.ErrorIs(t, err, ErrUserNotFound)
}
