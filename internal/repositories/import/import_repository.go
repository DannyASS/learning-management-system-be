package import_repository

import (
	"github.com/DannyAss/users/internal/database"
	import_model "github.com/DannyAss/users/internal/models/database_model/import"
	"gorm.io/gorm"
)

type importRepository struct {
	db *database.DBManager
	tx *gorm.DB
}

type IImportRepository interface {
	WithTx(tx *gorm.DB) IImportRepository

	InsertImportHDR(model import_model.ImportHdr) (*import_model.ImportHdr, error)
	UpdateImportHDR(model import_model.ImportHdr) error

	InsertImportDTL(model []import_model.ImportDtl) error

	GetImportByID(id int64) (*import_model.ImportHdr, error)
	GetImportByType(model import_model.ImportHdr) (*import_model.ImportHdr, error)
	GetImportByClassId(model import_model.ImportHdr) (string, error)
}

func NewImportRepository(db *database.DBManager) IImportRepository {
	return &importRepository{db: db}
}

func (r *importRepository) WithTx(tx *gorm.DB) IImportRepository {
	return &importRepository{db: r.db, tx: tx}
}

func (r *importRepository) getDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db.GetDB()
}

func (r *importRepository) InsertImportHDR(model import_model.ImportHdr) (*import_model.ImportHdr, error) {
	tx := r.getDB()

	if err := tx.Create(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *importRepository) UpdateImportHDR(model import_model.ImportHdr) error {
	tx := r.getDB()

	if err := tx.Model(&import_model.ImportHdr{}).
		Where("id = ?", model.ID).
		Updates(&model).Error; err != nil {
		return err
	}

	return nil
}

func (r *importRepository) InsertImportDTL(model []import_model.ImportDtl) error {
	tx := r.getDB()

	if len(model) == 0 {
		return nil
	}

	if err := tx.Create(&model).Error; err != nil {
		return err
	}

	return nil
}

func (r *importRepository) GetImportByID(id int64) (*import_model.ImportHdr, error) {
	tx := r.getDB()

	var data import_model.ImportHdr
	if err := tx.First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *importRepository) GetDetailImport(id int64) ([]import_model.ImportDtl, error) {
	var (
		data []import_model.ImportDtl
	)
	tx := r.getDB()

	if err := tx.Where("import_id = ?", id).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil

}

func (r *importRepository) GetImportByType(model import_model.ImportHdr) (*import_model.ImportHdr, error) {
	var (
		data import_model.ImportHdr
	)
	tx := r.getDB().Model(&data)

	if err := tx.Preload("Detail").Where("class_id = ?", model.ClassID).Where("type = ?", model.Type).Order("id desc").First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *importRepository) GetImportByClassId(model import_model.ImportHdr) (string, error) {
	var status string
	tx := r.getDB()

	tx = tx.Raw(`
			select status from import_hdr a
			where a.class_id = ?
			and a.type = ?
			order by id desc
			limit 1	
		`, model.ClassID, model.Type).Scan(&status)

	if err := tx.Error; err != nil {
		return "", err
	}

	return status, nil

}
