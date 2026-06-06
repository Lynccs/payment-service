package service

import (
	"context"

	"github.com/Lynccs/payment-service/internal/app/dto"
)

type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.UserResponse, error)
	GetProfile(ctx context.Context, userID int) (dto.UserResponse, error)
}
