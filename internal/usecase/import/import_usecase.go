package import_usecase

import (
	import_model "github.com/DannyAss/users/internal/models/database_model/import"
	import_repository "github.com/DannyAss/users/internal/repositories/import"
)

type importUsecase struct {
	repos import_repository.IImportRepository
}

type IImportUsecase interface {
	GetImport(typeImport string, id uint64) (*import_model.ImportHdr, error)
	// GetStatusImportByClassId(id uint64) (*map[string]string, error)
}

func NewImportUsecase(repo import_repository.IImportRepository) IImportUsecase {
	return &importUsecase{repos: repo}
}

func (i *importUsecase) GetImport(typeImport string, id uint64) (*import_model.ImportHdr, error) {
	model := import_model.ImportHdr{
		Type:    typeImport,
		ClassID: &id,
	}

	data, err := i.repos.GetImportByType(model)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (i *importUsecase) GetStatusImportByClassId(id uint64, types string) (*map[string]string, error) {
	status, err := i.repos.GetImportByClassId(import_model.ImportHdr{
		ClassID: &id,
		Type:    types,
	})

	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"status": status,
	}

	return &data, nil
}
