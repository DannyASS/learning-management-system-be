package absen_model

import "time"

type AbsenHdr struct {
	ID      uint64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Tanggal string  `json:"tanggal" gorm:"type:date;not null"`
	KelasID *uint64 `json:"kelas_id" gorm:"column:kelas_id"`
	MapelID *uint64 `json:"mapel_id" gorm:"column:mapel_id"`
	Catatan *string `json:"catatan" gorm:"type:text"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy *uint64    `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy *uint64    `json:"updated_by"`

	// relation ke detail
	Details []AbsenDtl `json:"details" gorm:"foreignKey:HdrID;constraint:OnDelete:CASCADE"`
}

// opsional: table name (kalau nama tabel bukan plural default)
func (AbsenHdr) TableName() string {
	return "absen_hdr"
}

type AbsenDtl struct {
	ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	HdrID     uint64 `json:"hdr_id" gorm:"not null"`
	StudentID uint64 `json:"student_id" gorm:"not null"`

	// ENUM: SAKIT, IZIN, ALFA
	Status     string  `json:"status" gorm:"type:enum('SAKIT','IZIN','ALFA');not null"`
	Keterangan *string `json:"keterangan" gorm:"type:text"`

	// relation ke header
	Header AbsenHdr `json:"header" gorm:"foreignKey:HdrID;references:ID"`
}

func (AbsenDtl) TableName() string {
	return "absen_dtl"
}

type DTOAbsenRequest struct {
	CourseId int           `json:"course_id"`
	Detail   []DTOAbsenDtl `json:"detail"`
}

type DTOAbsenDtl struct {
	StudentID  uint64  `json:"student_id" gorm:"not null"`
	Status     string  `json:"status" gorm:"type:enum('SAKIT','IZIN','ALFA');not null"`
	Keterangan *string `json:"keterangan" gorm:"type:text"`
}
