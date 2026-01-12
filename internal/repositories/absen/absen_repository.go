package absen_repository

import (
	"github.com/DannyAss/users/internal/database"
	absen_model "github.com/DannyAss/users/internal/models/database_model/absen"
	"gorm.io/gorm"
)

type absenRepository struct {
	db *database.DBManager
	tx *gorm.DB
}

type IAbsemRepository interface {
	WithTx(tx *gorm.DB) IAbsemRepository
	InsertAbsen(hdr absen_model.AbsenHdr, dtl []absen_model.AbsenDtl) (err error)
	GetAbsenToday(hdr absen_model.AbsenHdr) bool
}

func NewAbsenRepos(db *database.DBManager) IAbsemRepository {
	return &absenRepository{db: db}
}

func (a *absenRepository) WithTx(tx *gorm.DB) IAbsemRepository {
	return &absenRepository{db: a.db, tx: tx}
}

func (a *absenRepository) getDB() *gorm.DB {
	if a.tx != nil {
		return a.tx
	}

	return a.db.GetDB()
}

func (a *absenRepository) InsertAbsen(hdr absen_model.AbsenHdr, dtl []absen_model.AbsenDtl) (err error) {
	// mulai transaksi
	tx := a.getDB().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	if err = tx.Create(&hdr).Error; err != nil {
		return err
	}

	for i := range dtl {
		dtl[i].HdrID = hdr.ID
	}

	if len(dtl) > 0 {
		if err = tx.CreateInBatches(&dtl, 20).Error; err != nil {
			return err
		}
	}

	return nil
}

func (a *absenRepository) GetAbsenToday(hdr absen_model.AbsenHdr) bool {
	tx := a.getDB()

	if err := tx.Where(&hdr).First(&absen_model.AbsenHdr{}).Error; err != nil {
		return false
	}

	return true
}
