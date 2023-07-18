package repositories

import (
	"context"
	"fmt"

	customErr "github.com/niluwats/task-manager-auth-service/pkg/errors"
	"github.com/niluwats/task-manager-auth-service/pkg/models"
	"gorm.io/gorm"
)

type UserRepositoy interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, ID int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Insert(ctx context.Context, user models.User) (*models.User, error)
	DeleteByID(ctx context.Context, ID int) error
}

type UserRepoDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(dbClient *gorm.DB) UserRepoDB {
	return UserRepoDB{db: dbClient}
}

func (repo UserRepoDB) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	result := repo.db.Find(&users)
	if result.Error != nil {
		return nil, customErr.InternalError{Err: fmt.Sprint("Error querying all users")}
	}

	return users, nil
}

func (repo UserRepoDB) GetByID(ctx context.Context, ID int) (*models.User, error) {
	var user models.User
	result := repo.db.First(&user, ID)
	if result.Error != nil {
		return nil, customErr.NotFoundError{Err: fmt.Sprint("User not found")}
	}

	return &user, nil
}

func (repo UserRepoDB) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if result := repo.db.Where(&models.User{Email: user.Email}).First(&user); result.Error != nil {
		return nil, customErr.NotFoundError{Err: fmt.Sprint("Email not found")}
	}

	return &user, nil
}

func (repo UserRepoDB) Insert(ctx context.Context, user models.User) (*models.User, error) {
	retrievedUser, _ := repo.GetByEmail(ctx, user.Email)
	if retrievedUser != nil {
		return nil, customErr.ConflictError{Err: fmt.Sprint("User from this email already exists")}
	}

	result := repo.db.Create(&user)
	if result.Error != nil {
		return nil, customErr.InternalError{Err: fmt.Sprint("Error creating new user")}
	}

	fmt.Println(user.ID)
	fmt.Println(user)
	return &user, nil
}

func (repo UserRepoDB) DeleteByID(ctx context.Context, ID int) error {
	var user models.User
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
