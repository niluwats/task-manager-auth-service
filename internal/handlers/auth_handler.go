package handlers

import (
	"context"
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
	ViewUser(ctx context.Context, req *pb.ViewUserRequest) (*pb.ViewUserResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

type AuthHandlerImpl struct {
	pb.UnimplementedAuthServiceServer
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandlerImpl {
	return AuthHandlerImpl{service: service}
}

func (h AuthHandlerImpl) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ID, err := h.service.Register(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
		Status:   true,
	})

	if err != nil {
		return &pb.RegisterResponse{Status: getHttpCode(err), Message: err.Error()}, nil
	}

	return &pb.RegisterResponse{UserID: int64(ID), Status: http.StatusCreated, Message: "Register successful"}, nil
}

func (h AuthHandlerImpl) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ID, err := h.service.Login(ctx, domain.User{Email: req.Email, Password: req.Password})
	if err != nil {
		return &pb.LoginResponse{Status: getHttpCode(err), Message: err.Error()}, nil
	}

	token, err := utils.GenerateJWT(req.Email, ID)
	if err != nil {
		return &pb.LoginResponse{Status: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return &pb.LoginResponse{Token: token, Status: http.StatusOK, Message: "Login successful"}, nil
}

func (h AuthHandlerImpl) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := utils.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateTokenResponse{Status: http.StatusUnauthorized, Message: err.Error()}, nil
	}

	ID, err := h.service.GetUserIDByEmail(ctx, claims.Email)
	if err != nil {
		return &pb.ValidateTokenResponse{Status: http.StatusNotFound, Message: "User not found"}, nil
	}

	return &pb.ValidateTokenResponse{UserID: int64(ID), Status: http.StatusOK, Message: "Validate successful"}, nil
}

func (h AuthHandlerImpl) ViewUser(ctx context.Context, req *pb.ViewUserRequest) (*pb.ViewUserResponse, error) {
	user, err := h.service.GetUserByID(ctx, int(req.UserID))
	if err != nil {
		return &pb.ViewUserResponse{Status: getHttpCode(err), Message: err.Error()}, nil
	}

	return &pb.ViewUserResponse{UserID: int64(user.ID), Email: user.Email, Firstname: user.FirstName, Lastname: user.LastName, Activitystatus: user.Status, Status: http.StatusOK, Message: "Success"}, nil
}

func getHttpCode(err error) int32 {
	switch err.(type) {
	case *errors.BadRequest:
		return http.StatusBadRequest
	case *errors.Unauthorized:
		return http.StatusUnauthorized
	case *errors.ConflictError:
		return http.StatusConflict
	case *errors.NotFoundError:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
