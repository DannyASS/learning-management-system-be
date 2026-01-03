package courses_model

import (
	"time"

	user_model "github.com/DannyAss/users/internal/models/database_model/user"
	"gorm.io/gorm"
)

// Course represents the courses table
type Course struct {
	ID           uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title        string  `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description  *string `gorm:"column:description;type:text" json:"description,omitempty"`
	ThumbnailURL *string `gorm:"column:thumbnail_url;type:varchar(255)" json:"thumbnail_url,omitempty"`
	CategoryID   *uint   `gorm:"column:category_id" json:"category_id,omitempty"`

	// status enum('draft','published','archived') DEFAULT 'draft'
	Status string `gorm:"column:status;type:enum('draft','published','archived');default:'draft'" json:"status"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy string         `gorm:"column:created_by;type:varchar(255)" json:"created_by,omitempty"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy *string        `gorm:"column:updated_by;type:varchar(255)" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// TableName explicit (optional)
func (Course) TableName() string {
	return "courses"
}

type FilterCourse struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	CategoryID int    `json:"category_id"`
	Status     string `json:"status"`
}

type CourseDTO struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnail_url"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	Status       string `json:"status"`
}

type CourseRequest struct {
	Page      int    `query:"page"`
	Perpage   int    `query:"perpage"`
	TotalData int    `query:"total_data"`
	TotalPage int    `query:"total_page"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortType  string `query:"sort_type"`
}

type Category struct {
	Id          uint   `gorm:"column:id"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

type AvailableCourse struct {
	Teacher []user_model.User `json:"teacher"`
	Courses []Course          `json:"courses"`
}
