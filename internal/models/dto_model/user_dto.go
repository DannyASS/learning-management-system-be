package dtomodel

import (
	"gorm.io/gorm"
)

type UserDTO struct {
	ID        uint64
	Name      string
	Phone     string
	DeletedAt gorm.DeletedAt
}
