package modules_usecase

import (
	"time"

	"github.com/DannyAss/users/internal/database"
	courses_model "github.com/DannyAss/users/internal/models/database_model/courses"
	modules_model "github.com/DannyAss/users/internal/models/database_model/modules"
	courses_repository "github.com/DannyAss/users/internal/repositories/courses"
	modules_repository "github.com/DannyAss/users/internal/repositories/modules"
)

type moduleUsecase struct {
	db          *database.DBManager
	moduleRepos modules_repository.IModulesRepos
	courseRepos courses_repository.ICourseRepos
}

type IModuleUsecase interface {
	GetListModule(req modules_model.ModuleRequest) (*map[string]interface{}, error)
	CreateModule(req modules_model.ModuleDTO) error
	GetModule(req modules_model.FilterModule) (*modules_model.ModuleDTO, error)
	UpdateModule(filter modules_model.FilterModule, data modules_model.ModuleDTO) error
	DeleteModule(filter modules_model.FilterModule) error
	GetCourse() ([]map[string]interface{}, error)
}

func NewModuleUsecase(db *database.DBManager, repos modules_repository.IModulesRepos, crepos courses_repository.ICourseRepos) IModuleUsecase {
	return &moduleUsecase{
		db:          db,
		moduleRepos: repos,
		courseRepos: crepos,
	}
}

func (m *moduleUsecase) GetListModule(req modules_model.ModuleRequest) (*map[string]interface{}, error) {
	var dtolist []modules_model.ModuleDTO

	filter := modules_model.FilterModule{
		Title: req.Search,
	}

	getdata, page, err := m.moduleRepos.GetList(filter, "", &req)
	if err != nil {
		return nil, err
	}

	for _, c := range getdata {
		dto := modules_model.ModuleDTO{
			ID:       c.ID,
			Title:    c.Title,
			CourseID: c.CourseID,
			Type:     c.Type,
		}

		dtolist = append(dtolist, dto)
	}

	response := map[string]interface{}{
		"data":       dtolist,
		"page":       page.Page,
		"perpage":    page.Perpage,
		"total_data": page.TotalData,
		"total_page": page.TotalPage,
	}

	return &response, nil
}

func (m *moduleUsecase) CreateModule(req modules_model.ModuleDTO) error {
	model := modules_model.Module{
		Title:      req.Title,
		CourseID:   req.CourseID,
		Type:       req.Type,
		OrderIndex: req.OrderIndex,
		CreatedAt:  time.Now(),
	}

	_, err := m.moduleRepos.InsertModule(model)
	if err != nil {
		return err
	}
	return nil
}

func (m *moduleUsecase) GetModule(req modules_model.FilterModule) (*modules_model.ModuleDTO, error) {

	data, err := m.moduleRepos.GetModule(req, "")
	if err != nil {
		return nil, err
	}

	resonse := modules_model.ModuleDTO{
		Title:      data.Title,
		CourseID:   data.CourseID,
		Type:       data.Type,
		OrderIndex: data.OrderIndex,
	}

	return &resonse, nil
}

func (m *moduleUsecase) GetCourse() ([]map[string]interface{}, error) {
	var data []map[string]interface{}

	course, _, err := m.courseRepos.Getlist(courses_model.FilterCourse{Status: "published"}, "", nil)
	if err != nil {
		return nil, err
	}

	for _, c := range course {
		dto := map[string]interface{}{
			"course_title": c.Title,
			"course_id":    c.ID,
		}

		data = append(data, dto)
	}

	return data, nil
}

func (m *moduleUsecase) UpdateModule(filter modules_model.FilterModule, data modules_model.ModuleDTO) error {
	model := modules_model.Module{
		Title:      data.Title,
		CourseID:   data.CourseID,
		Type:       data.Type,
		OrderIndex: data.OrderIndex,
	}
	err := m.moduleRepos.UpdateData(filter, model)
	if err != nil {
		return err
	}
	return nil
}

func (m *moduleUsecase) DeleteModule(filter modules_model.FilterModule) error {
	return nil
}
