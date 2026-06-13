package services

import (
	"context"
	"fmt"
	"unicode"

	"github.com/Lynccs/payment-service/internal/app/dto"
	"github.com/Lynccs/payment-service/internal/app/models"
	"github.com/Lynccs/payment-service/internal/app/repo"
	appsvc "github.com/Lynccs/payment-service/internal/app/service"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists  = fmt.Errorf("email already exists")
	ErrPasswordTooShort    = fmt.Errorf("password must be at least 8 characters")
	ErrPasswordNoUpper     = fmt.Errorf("password must contain at least one uppercase letter")
	ErrPasswordNoDigit     = fmt.Errorf("password must contain at least one digit")
	ErrPasswordNoSymbol    = fmt.Errorf("password must contain at least one special character")
	ErrInvalidCredentials  = fmt.Errorf("invalid email or password")
	ErrUserNotFound        = fmt.Errorf("user not found")
)

type userService struct {
	userRepo   repo.UserRepository
	walletRepo repo.WalletRepository
}

func NewUserService(userRepo repo.UserRepository, walletRepo repo.WalletRepository) appsvc.UserService {
	return &userService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

func (s *userService) Register(ctx context.Context, req dto.RegisterRequest) (dto.UserResponse, error) {
	op := "userService.Register"
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return dto.UserResponse{}, ErrEmailAlreadyExists
	}

	if err := validatePassword(req.Password); err != nil {
		return dto.UserResponse{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("%s: hash password: %w", op, err)
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}

	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("%s: create user: %w", op, err)
	}

	if _, err := s.walletRepo.Create(ctx, userID); err != nil {
		return dto.UserResponse{}, fmt.Errorf("%s: create wallet: %w", op, err)
	}

	createdUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("%s: get created user: %w", op, err)
	}

	return dto.UserResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

func (s *userService) Login(ctx context.Context, req dto.LoginRequest) (dto.UserResponse, error) {
	op := "userService.Login"
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.UserResponse{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	var hasUpper, hasDigit, hasSymbol bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSymbol = true
		}
	}
	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasDigit {
		return ErrPasswordNoDigit
	}
	if !hasSymbol {
		return ErrPasswordNoSymbol
	}
	return nil
}

func (s *userService) GetProfile(ctx context.Context, userID int) (dto.UserResponse, error) {
	op := "userService.GetProfile"
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
