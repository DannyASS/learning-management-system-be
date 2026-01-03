package menus_model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// StringArray is a helper type to store []string as JSON in MySQL.
type StringArray []string

// Value implements driver.Valuer — called when saving to DB.
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		// store empty JSON array instead of NULL
		return "[]", nil
	}
	b, err := json.Marshal([]string(s))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements sql.Scanner — called when reading from DB.
func (s *StringArray) Scan(src interface{}) error {
	if src == nil {
		*s = nil
		return nil
	}

	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return errors.New("incompatible type for StringArray")
	}

	// if empty string, set nil or empty slice
	if len(data) == 0 {
		*s = nil
		return nil
	}

	var dst []string
	if err := json.Unmarshal(data, &dst); err != nil {
		return err
	}
	*s = dst
	return nil
}

// ------------------ Menu model ------------------
type Menu struct {
	ID         uint        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Path       string      `gorm:"column:path;type:varchar(255);not null" json:"path"`
	URL        string      `gorm:"column:url;type:varchar(255);not null" json:"url"`
	Permission StringArray `gorm:"column:permission;type:json;not null" json:"permission"`
	Name       string      `gorm:"column:name;type:varchar(100);not null" json:"name"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy string         `gorm:"column:created_by;type:varchar(255);default:null" json:"created_by"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy *string        `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName explicit (optional)
func (Menu) TableName() string {
	return "menus"
}

// ------------------ RoleMenu model ------------------
type RoleMenu struct {
	ID         uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID     uint           `gorm:"column:role_id;not null" json:"role_id"`
	Permission *string        `gorm:"column:permission;type:text;default:null" json:"permission"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  *time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (RoleMenu) TableName() string {
	return "role_menus"
}
