package user_model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	Username   string         `gorm:"type:varchar(100);not null" json:"username"`
	Email      string         `gorm:"type:Text;not null" json:"email"`
	Password   string         `gorm:"type:Text;not null" json:"-"`
	Gender     string         `gorm:"type:enum('male', 'female');not null" json:"-"`
	Phone      string         `gorm:"type:Text" json:"phone,omitempty"`
	IsVerified bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type FilterUser struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	IsVerified bool   `json:"is_verified"`
}

type UserRole struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	UserID    uint64         `gorm:"not null"`
	RoleID    uint64         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // untuk soft delete
}

// Pastikan GORM pakai tabel 'user_roles'
func (UserRole) TableName() string {
	return "user_roles"
}

type UserRoleDTO struct {
	UserId uint64 `json:"user_id"`
	RoleId uint64 `json:"role_id"`
}

type RegisterDTO struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserToken struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint
	RefreshToken string `gorm:"type:text"`
	UserAgent    string
	IPAddress    string
	ExpiresAt    time.Time
	Revoked      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (UserToken) TableName() string {
	return "user_tokens"
}

type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type FilterUserToken struct {
	RefreshToken string `json:"refresh_token"`
	Revoked      bool   `json:"revoked"`
}

type UserDTO struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     *Role  `json:"role"`
}

type Pagination struct {
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortType  string `query:"sort_type"`
	Page      int    `query:"page"`
	Perpage   int    `query:"perpage"`
	TotalData int    `query:"total_data"`
	TotalPage int    `query:"total_page"`
}

type RequestUserGetList struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role_name" gorm:"column:role_name"`
	RoleId   int    `json:"role_id" gorm:"column:role_id"`
}
