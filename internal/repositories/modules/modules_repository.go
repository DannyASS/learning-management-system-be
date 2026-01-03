package modules_repository

import (
	"math"
	"strings"

	"github.com/DannyAss/users/internal/database"
	modules_model "github.com/DannyAss/users/internal/models/database_model/modules"
	"gorm.io/gorm"
)

type modulesRepos struct {
	db *database.DBManager
	tx *gorm.DB
}

type IModulesRepos interface {
	WithTx(tx *gorm.DB) IModulesRepos
	InsertModule(model modules_model.Module) (*modules_model.Module, error)
	UpdateData(filter modules_model.FilterModule, model modules_model.Module) error
	GetList(filter modules_model.FilterModule, mode string, page *modules_model.ModuleRequest) ([]modules_model.Module, *modules_model.ModuleRequest, error)
	GetModule(filter modules_model.FilterModule, mode string) (*modules_model.Module, error)
	DeleteModule(filter modules_model.FilterModule) error
}

func NewModulesRepos(db *database.DBManager) IModulesRepos {
	return &modulesRepos{db: db}
}

func (m *modulesRepos) WithTx(tx *gorm.DB) IModulesRepos {
	return &modulesRepos{db: m.db, tx: tx}
}

func (m *modulesRepos) getDB() *gorm.DB {
	if m.tx != nil {
		return m.tx
	}
	return m.db.GetDB()
}

func (m *modulesRepos) InsertModule(model modules_model.Module) (*modules_model.Module, error) {
	tx := m.getDB()

	if err := tx.Create(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (m *modulesRepos) UpdateData(filter modules_model.FilterModule, model modules_model.Module) error {
	tx := m.getDB()

	if err := tx.Where(&filter).Updates(&model).Error; err != nil {
		return err
	}

	return nil
}

func (m *modulesRepos) GetList(filter modules_model.FilterModule, mode string, page *modules_model.ModuleRequest) ([]modules_model.Module, *modules_model.ModuleRequest, error) {
	tx := m.getDB().Debug().Model(&modules_model.Module{})

	var (
		modules   []modules_model.Module
		totalData int64
	)

	// kalau pakai pagination otomatis pakai mode LIKE
	if page != nil {
		mode = "like"
	}

	switch mode {
	case "full":
		tx = tx.Where(filter)

	case "like":
		name := filter.Title
		filter.Title = ""

		if name != "" {
			tx = tx.Where("name LIKE ?", "%"+name+"%").
				Where(&filter)
		}

		// Hitung total
		if err := tx.Count(&totalData).Error; err != nil {
			return nil, nil, err
		}

		// Pagination
		if err := tx.Order(page.SortBy + " " + page.SortType).
			Limit(page.Perpage).
			Offset((page.Page - 1) * page.Perpage).
			Find(&modules).Error; err != nil {
			return nil, nil, err
		}

		page.TotalPage = int(math.Ceil(float64(totalData) / float64(page.Perpage)))
		page.TotalData = int(totalData)

		return modules, page, nil

	default:
		tx = tx.Where(&filter)
	}

	if err := tx.Find(&modules).Error; err != nil {
		return nil, nil, err
	}

	return modules, nil, nil
}

func (m *modulesRepos) GetModule(filter modules_model.FilterModule, mode string) (*modules_model.Module, error) {
	tx := m.getDB()

	var module modules_model.Module

	switch strings.ToLower(mode) {
	case "full":
		tx = tx.Where(filter)
	default:
		tx = tx.Where(filter)
	}

	if err := tx.First(&module).Error; err != nil {
		return nil, err
	}

	return &module, nil
}

func (m *modulesRepos) DeleteModule(filter modules_model.FilterModule) error {
	tx := m.getDB()

	if err := tx.Where(&filter).Delete(&modules_model.Module{}).Error; err != nil {
		return err
	}

	return nil
}
