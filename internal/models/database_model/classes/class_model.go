package classes_model

import (
	"time"

	import_model "github.com/DannyAss/users/internal/models/database_model/import"
	"gorm.io/gorm"
)

// ================================
// class_hdr
// ================================
type ClassHdr struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `json:"name"`
	TeacherID uint   `json:"teacher_id"`
	Status    uint8  `json:"status"`
	Level     string `json:"level"`
	Major     string `json:"major"`

	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`

	Modules    []ClassModule          `gorm:"foreignKey:ClassID" json:"modules"`
	Students   []ClassStudent         `gorm:"foreignKey:ClassID" json:"students"`
	ImportHdrs import_model.ImportHdr `gorm:"foreignKey:ClassID" json:"imports"`
}

func (ClassHdr) TableName() string {
	return "class_hdr"
}

// ================================
// course_modules (reference)
// ================================
type CourseModule struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID   uint   `json:"course_id"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	OrderIndex int    `json:"order_index"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func (CourseModule) TableName() string {
	return "course_modul"
}

// ================================
// class_module
// ================================
type ClassModule struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassID    uint64 `json:"class_id"`
	ModuleID   uint   `json:"module_id"`
	ModuleName string `json:"module_name"`
	Status     string `json:"status"`

	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`

	ClassHdr     *ClassHdr     `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	CourseModule *CourseModule `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
}

func (ClassModule) TableName() string {
	return "class_module"
}

// ================================
// class_student
// ================================
type ClassStudent struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassID     uint64 `json:"class_id"`
	UserID      uint64 `json:"user_id"`
	Gender      string `json:"gender"`
	Status      string `json:"status"`
	StudentName string `gorm:"colomn:student_name" json:"student_name"`

	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *uint      `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *uint      `json:"updated_by"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`

	ClassHdr *ClassHdr `gorm:"foreignKey:ClassID" json:"class,omitempty"`
}

func (ClassStudent) TableName() string {
	return "class_student"
}

// ================================
// class_assignment
// ================================
type ClassAssignment struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassID     uint64     `json:"class_id"`
	ModuleID    uint       `json:"module_id"`
	Title       string     `json:"title"`
	LinkURL     string     `json:"link_url"`
	Description string     `json:"description"`
	ExpiredAt   *time.Time `json:"expired_at"`
	Status      string     `json:"status"`

	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *uint      `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *uint      `json:"updated_by"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`

	ClassHdr     *ClassHdr     `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	CourseModule *CourseModule `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
}

func (ClassAssignment) TableName() string {
	return "class_assigment"
}

type ClassHDRRequest struct {
	Page      int    `query:"page"`
	Perpage   int    `query:"perpage"`
	TotalData int    `query:"total_data"`
	TotalPage int    `query:"total_page"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortType  string `query:"sort_type"`
}

type ClassHDRQuery struct {
	Id              int    `gorm:"column:id"`
	Name            string `gorm:"column:name"`
	CompletedModule int    `gorm:"column:completed_module"`
	TotalModule     int    `gorm:"column:total_module"`
	TotalStudent    int    `gorm:"column:total_student"`
}

type ClassHDRListDTOL struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Progress     float64 `json:"progress"`
	TotalStudent int     `json:"total_student"`
}

type ClassHdrDTO struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `json:"name"`
	TeacherID uint   `json:"teacher_id"`
	Status    uint8  `json:"status"`
	Level     string `json:"level"`
	Major     string `json:"major"`
}

type StudentAvailable struct {
	UserID uint64 `json:"user_id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type ClassCourse struct {
	ID        uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	ClassID   uint64 `gorm:"column:class_id;not null;index"`
	CourseID  uint64 `gorm:"column:course_id;not null;index"`
	TeacherID uint64 `gorm:"column:teacher_id;not null;index"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	CreatedBy *string        `gorm:"column:created_by;size:255"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy *string        `gorm:"column:updated_by;size:255"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ClassCourse) TableName() string {
	return "class_course"
}

type DashBoardInformation struct {
	TotalCourses        int64 `json:"total_course"`
	ActiveAssigment     int64 `json:"active_assigment"`
	TotalCompliteCourse int64 `json:"total_compliteCourse"`
}
