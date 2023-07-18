package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/niluwats/task-manager-auth-service/pkg/models"
	"github.com/niluwats/task-manager-auth-service/pkg/pb"
	"github.com/niluwats/task-manager-auth-service/pkg/repositories"
	"github.com/niluwats/task-manager-auth-service/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *pb.RegisterRequest) *pb.RegisterResponse
	Login(ctx context.Context, req *pb.LoginRequest) *pb.LoginResponse
	Validate(ctx context.Context, req *pb.ValidateTokenRequest) *pb.ValidateTokenResponse
}

type DefaultUserService struct {
	repo repositories.UserRepositoy
	Jwt  utils.JWTWrapper
}

func NewUserService(repo repositories.UserRepositoy, wrapper utils.JWTWrapper) DefaultUserService {
	return DefaultUserService{repo, wrapper}
}

func (service DefaultUserService) Register(ctx context.Context, req *pb.RegisterRequest) *pb.RegisterResponse {
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPw),
	}

	newUser, err := service.repo.Insert(ctx, user)
	if err != nil {
		return &pb.RegisterResponse{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return &pb.RegisterResponse{Status: http.StatusCreated, Message: fmt.Sprintf("%v", newUser.ID)}
}

func (service DefaultUserService) Login(ctx context.Context, req *pb.LoginRequest) *pb.LoginResponse {
	user, err := service.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return &pb.LoginResponse{Status: http.StatusNotFound, Message: err.Error()}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &pb.LoginResponse{Status: http.StatusUnauthorized, Message: "Incorrect password"}
	}

	token, err := service.Jwt.GenerateJWT(req.Email, int(user.ID))

	return &pb.LoginResponse{Token: token, Status: http.StatusOK, Message: "Login success"}
}
