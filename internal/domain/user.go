package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email,omitempty" gorm:"not null"`
	Password  string `json:"password,omitempty" gorm:"not null"`
	Status    bool   `json:"status"`
}
