package service

import (
	"context"

	"github.com/niluwats/task-manager-auth-service/internal/domain"
	customErr "github.com/niluwats/task-manager-auth-service/internal/errors"
	"github.com/niluwats/task-manager-auth-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, user domain.User) (uint, error)
	Login(ctx context.Context, user domain.User) (uint, error)
	GetUserIDByEmail(ctx context.Context, email string) (uint, error)
	GetUserByID(ctx context.Context, ID int) (*domain.User, error)
}

type DefaultUserService struct {
	repo repositories.UserRepositoy
}

func NewUserService(repo repositories.UserRepositoy) DefaultUserService {
	return DefaultUserService{repo}
}

func (service DefaultUserService) Register(ctx context.Context, newUser domain.User) (uint, error) {
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	newUser.Password = string(hashedPw)

	userID, err := service.repo.Insert(ctx, newUser)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (service DefaultUserService) Login(ctx context.Context, loggedUser domain.User) (uint, error) {
	user, err := service.repo.GetByEmail(ctx, loggedUser.Email)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loggedUser.Password))

	if err != nil {
		return 0, &customErr.Unauthorized{Err: "Incorrect password"}
	}
	return user.ID, nil
}

func (service DefaultUserService) GetUserIDByEmail(ctx context.Context, email string) (uint, error) {
	user, err := service.repo.GetByEmail(ctx, email)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (service DefaultUserService) GetUserByID(ctx context.Context, ID int) (*domain.User, error) {
	user, err := service.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
