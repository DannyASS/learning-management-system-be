package user_model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type RoleFilter struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type UserRoleTokenDTO struct {
	Id       uint
	Username string
	Email    string
	RolesId  []uint
}
