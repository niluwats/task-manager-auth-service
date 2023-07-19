package repositories

import (
	"context"
	"fmt"

	"github.com/niluwats/task-manager-auth-service/internal/domain"
	customErr "github.com/niluwats/task-manager-auth-service/internal/errors"
	"gorm.io/gorm"
)

type UserRepositoy interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, ID int) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Insert(ctx context.Context, user domain.User) (uint, error)
	DeleteByID(ctx context.Context, ID int) error
}

type UserRepoDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(dbClient *gorm.DB) UserRepoDB {
	return UserRepoDB{db: dbClient}
}

func (repo UserRepoDB) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	result := repo.db.Find(&users)
	if result.Error != nil {
		return nil, customErr.InternalError{Err: fmt.Sprint("Error querying all users")}
	}

	return users, nil
}

func (repo UserRepoDB) GetByID(ctx context.Context, ID int) (*domain.User, error) {
	var user domain.User
	result := repo.db.First(&user, ID)
	if result.Error != nil {
		return nil, customErr.NotFoundError{Err: fmt.Sprint("User not found")}
	}

	return &user, nil
}

func (repo UserRepoDB) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if result := repo.db.Where(&domain.User{Email: user.Email}).First(&user); result.Error != nil {
		return nil, customErr.NotFoundError{Err: fmt.Sprint("Email not found")}
	}

	return &user, nil
}

func (repo UserRepoDB) Insert(ctx context.Context, user domain.User) (uint, error) {
	retrievedUser, _ := repo.GetByEmail(ctx, user.Email)
	if retrievedUser != nil {
		return 0, customErr.ConflictError{Err: fmt.Sprint("User from this email already exists")}
	}

	result := repo.db.Create(&user)
	if result.Error != nil {
		return 0, customErr.InternalError{Err: fmt.Sprint("Error creating new user")}
	}

	fmt.Println(user.ID)
	fmt.Println(user)
	return user.ID, nil
}

func (repo UserRepoDB) DeleteByID(ctx context.Context, ID int) error {
	var user domain.User
	user.ID = uint(ID)
	result := repo.db.Model(&user).Update("status", false)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return customErr.NotFoundError{Err: fmt.Sprint("User not found")}
		}
		return customErr.InternalError{Err: fmt.Sprint("Error updaring user status")}
	}

	return nil
}
