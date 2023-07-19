package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/niluwats/task-manager-auth-service/api/pb"
	"github.com/niluwats/task-manager-auth-service/internal/domain"
	"github.com/niluwats/task-manager-auth-service/internal/errors"
	"github.com/niluwats/task-manager-auth-service/internal/service"
	"github.com/niluwats/task-manager-auth-service/internal/utils"
)

type AuthHandlers interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

type DefaultAuthHandler struct {
	pb.UnimplementedAuthServiceServer
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) DefaultAuthHandler {
	return DefaultAuthHandler{service: service}
}

func (h DefaultAuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ID, err := h.service.Register(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
		Status:   true,
	})

	if err != nil {
		return &pb.RegisterResponse{Status: getHttpCode(err), Message: err.Error()}, err
	}

	return &pb.RegisterResponse{Status: http.StatusCreated, Message: fmt.Sprintf("%v", ID)}, nil
}

func (h DefaultAuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ID, err := h.service.Login(ctx, domain.User{Email: req.Email, Password: req.Password})
	if err != nil {
		return &pb.LoginResponse{Status: getHttpCode(err), Message: err.Error()}, err
	}

	token, err := utils.GenerateJWT(req.Email, ID)
	if err != nil {
		return &pb.LoginResponse{Status: http.StatusInternalServerError, Message: err.Error()}, err
	}

	return &pb.LoginResponse{Token: token, Status: http.StatusOK, Message: "Login Success"}, nil
}

func (h DefaultAuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := utils.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateTokenResponse{Status: http.StatusUnauthorized, Message: err.Error()}, err
	}

	ID, err := h.service.GetUserIDByEmail(ctx, claims.Email)
	if err != nil {
		return &pb.ValidateTokenResponse{Status: http.StatusNotFound, Message: "User not found"}, err
	}

	return &pb.ValidateTokenResponse{UserID: int32(ID), Status: http.StatusOK, Message: "Validate successful"}, nil
}

func getHttpCode(err error) int32 {
	switch err {
	case errors.BadRequest{}:
		return http.StatusBadRequest
	case errors.Unauthorized{}:
		return http.StatusUnauthorized
	case errors.ConflictError{}:
		return http.StatusConflict
	case errors.NotFoundError{}:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
