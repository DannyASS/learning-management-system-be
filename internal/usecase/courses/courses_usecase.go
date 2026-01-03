package courses_usecase

import (
	"errors"
	"time"

	"github.com/DannyAss/users/internal/database"
	courses_model "github.com/DannyAss/users/internal/models/database_model/courses"
	courses_repository "github.com/DannyAss/users/internal/repositories/courses"
	"github.com/DannyAss/users/pkg/utils"
)

type courseUsecase struct {
	db          *database.DBManager
	courseRepos courses_repository.ICourseRepos
}

type ICourseUsecase interface {
	GetListCourse(req courses_model.CourseRequest) (*map[string]interface{}, error)
	CreateCourse(req courses_model.CourseDTO) error
	GetCourse(req courses_model.FilterCourse) (*courses_model.CourseDTO, error)
	UpdateCourse(req courses_model.FilterCourse, data courses_model.CourseDTO) error
}

func NewCourseUsecase(db *database.DBManager, repos courses_repository.ICourseRepos) ICourseUsecase {
	return &courseUsecase{db: db, courseRepos: repos}
}

func (c *courseUsecase) GetListCourse(req courses_model.CourseRequest) (*map[string]interface{}, error) {
	allowSort := map[string]bool{
		"title":  true,
		"id":     true,
		"status": true,
	}

	filter := courses_model.FilterCourse{
		Title: req.Search,
	}

	if req.SortType != "asc" && req.SortType != "desc" {
		req.SortType = "asc"
	}

	if _, ok := allowSort[req.SortBy]; !ok {
		req.SortBy = "id"
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Perpage == 0 {
		req.Perpage = 10
	}

	var courses []courses_model.CourseDTO

	dataDB, page, err1 := c.courseRepos.Getlist(filter, "", &req)
	if err1 != nil {
		return nil, err1
	}

	categories, err2 := c.courseRepos.GetCategory()
	if err2 != nil {
		return nil, err2
	}

	for _, c := range dataDB {
		nameCategories := utils.FilterSlice(categories, func(k courses_model.Category) bool {
			if c.CategoryID == nil {
				return false
			}
			return k.Id == uint(*c.CategoryID)
		})

		dto := courses_model.CourseDTO{
			Title:        c.Title,
			Id:           int(c.ID),
			Description:  *c.Description,
			Status:       c.Status,
			CategoryID:   int(*c.CategoryID),
			CategoryName: nameCategories[0].Name,
		}

		courses = append(courses, dto)
	}

	response := map[string]interface{}{
		"page":       page.Page,
		"perpage":    page.Perpage,
		"total_data": page.TotalData,
		"total_page": page.TotalPage,
		"data":       courses,
	}

	return &response, nil
}

func (c *courseUsecase) CreateCourse(req courses_model.CourseDTO) error {
	model := courses_model.Course{
		Title:        req.Title,
		Description:  &req.Description,
		ThumbnailURL: &req.ThumbnailURL,
		CategoryID:   utils.UIntPtr(req.CategoryID),
		Status:       "draft",
		CreatedAt:    time.Now(),
		CreatedBy:    "system",
	}

	if _, err := c.courseRepos.InsertCourse(model); err != nil {
		return err
	}

	return nil
}

func (c *courseUsecase) GetCourse(req courses_model.FilterCourse) (*courses_model.CourseDTO, error) {
	data, err := c.courseRepos.GetCourse(req, "")
	if err != nil {
		return nil, err
	}

	id := data.CategoryID
	categories, err2 := c.courseRepos.GetCategoryById(int(*id))
	if err2 != nil {
		return nil, err2
	}

	response := courses_model.CourseDTO{
		Title:        data.Title,
		Id:           int(data.ID),
		Description:  *data.Description,
		Status:       data.Status,
		CategoryID:   int(*data.CategoryID),
		CategoryName: categories.Name,
	}

	return &response, nil
}

func (c *courseUsecase) UpdateCourse(req courses_model.FilterCourse, data courses_model.CourseDTO) error {
	timeNow := time.Now()

	model := courses_model.Course{
		Title:        data.Title,
		Description:  &data.Description,
		ThumbnailURL: &data.ThumbnailURL,
		Status:       data.Status,
		UpdatedAt:    &timeNow,
	}

	cekStatus, errCek := c.courseRepos.GetCourse(req, "")
	if errCek != nil {
		return errCek
	}

	if cekStatus.Status != "draft" && data.Status == "draft" {
		return errors.New("Cannot update status")
	}

	if cekStatus.Status != "draft" && data.Title != cekStatus.Title {
		return errors.New("You cannot update the title because this module is no longer in draft status.")
	}

	if cekStatus.Status != "draft" && data.Description != *cekStatus.Description {
		return errors.New("You cannot update the description because this module is no longer in draft status.")
	}

	if err := c.courseRepos.UpdateData(req, model); err != nil {
		return err
	}

	return nil
}
