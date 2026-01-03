package import_model

import "time"

type ImportHdr struct {
	ID        int64     `json:"id" db:"id"`
	FileName  string    `json:"file_name" db:"file_name"`
	Type      string    `json:"type" db:"type"`
	OnStep    *string   `json:"on_step" db:"on_step"`
	ClassID   *uint64   `gorm:"column:class_id" json:"class_id,omitempty"`
	Desc      string    `json:"desc" db:"desc"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	Detail []ImportDtl `gorm:"foreignKey:ImportID" json:"details"`
}

func (ImportHdr) TableName() string {
	return "import_hdr"
}

type ImportDtl struct {
	ID        int64   `json:"id" db:"id"`
	ImportID  int64   `json:"import_id" db:"import_id"`
	FieldName *string `json:"field_name" db:"field_name"` // nullable
	Desc      *string `json:"desc" db:"desc"`             // nullable
}

func (ImportDtl) TableName() string {
	return "import_dtl"
}
