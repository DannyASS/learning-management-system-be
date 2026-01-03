package modules_model

import (
	"time"

	"gorm.io/gorm"
)

// Module represents the modules table
type Module struct {
	ID         uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CourseID   uint   `gorm:"column:course_id;not null" json:"course_id"`
	Title      string `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Type       string `gorm:"column:type;type:enum('video','article','quiz','assignment');not null" json:"type"`
	OrderIndex int    `gorm:"column:order_index;default:1" json:"order_index"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// TableName override (optional)
func (Module) TableName() string {
	return "course_modules"
}

// Filter struct for searching/filtering modules (optional)
type FilterModule struct {
	ID       int    `json:"id"`
	CourseID int    `json:"course_id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

// DTO struct for returning module + course name (optional)
type ModuleDTO struct {
	ID         uint   `json:"id"`
	CourseID   uint   `json:"course_id"`
	CourseName string `json:"course_name"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	OrderIndex int    `json:"order_index"`
}

// Query request for pagination/filtering
type ModuleRequest struct {
	Page      int    `query:"page"`
	Perpage   int    `query:"perpage"`
	TotalData int    `query:"total_data"`
	TotalPage int    `query:"total_page"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortType  string `query:"sort_type"`
}
