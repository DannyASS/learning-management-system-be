package classes_usecase

import (
	"errors"
	"fmt"

	"github.com/DannyAss/users/internal/database"
	classes_model "github.com/DannyAss/users/internal/models/database_model/classes"
	courses_model "github.com/DannyAss/users/internal/models/database_model/courses"
	classes_repository "github.com/DannyAss/users/internal/repositories/classes"
	courses_repository "github.com/DannyAss/users/internal/repositories/courses"
	import_repository "github.com/DannyAss/users/internal/repositories/import"
	jobs "github.com/DannyAss/users/internal/worker/job"
	job_excel "github.com/DannyAss/users/internal/worker/job/excel"
	"github.com/DannyAss/users/pkg/utils"
)

type classUsecase struct {
	repo  classes_repository.IClassesRepository
	irepo import_repository.IImportRepository
	crepo courses_repository.ICourseRepos
	db    *database.DBManager
}

var (
	InternalServerError = errors.New("Internal Server Error")
	UnprocessableEntity = errors.New("Unprocessable Entity")
)

type IClassesUsecase interface {
	GetListClassesPage(req classes_model.ClassHDRRequest, teacherId int) (*map[string]interface{}, error)
	CreateClass(req classes_model.ClassHdrDTO) error
	UpdateClass(req classes_model.ClassHdrDTO, StudentFilePath string, CoursesFilePath string) error
	GetAvailableStudent() ([]classes_model.StudentAvailable, error)
	GetClassByID(id uint64) (*classes_model.ClassHdr, error)
	DeleteStudentClassByID(id uint64) error
	DeleteStudentClassByIDClass(id uint64) error
	GetStudentClassByIDClass(id uint64) ([]classes_model.ClassStudent, error)
	GetAvailableCourseCLass() (*courses_model.AvailableCourse, error)
	GetAllCourseClass(id uint64) ([]map[string]interface{}, error)
	GetInformDashboardClass(id int, teacherId int) (map[string]interface{}, error)
	GetAllModulByClassAndRole(classId int, teacherId int) ([]map[string]interface{}, error)
	GetAvailableModulDash(classId int, teacherId int) ([]map[string]interface{}, error)
}

func NewClassesUsecase(repo classes_repository.IClassesRepository, irepo import_repository.IImportRepository, db *database.DBManager, crepo courses_repository.ICourseRepos) IClassesUsecase {
	return &classUsecase{repo: repo, irepo: irepo, db: db, crepo: crepo}
}

func (c *classUsecase) GetListClassesPage(req classes_model.ClassHDRRequest, teacherId int) (*map[string]interface{}, error) {
	var list []classes_model.ClassHDRListDTOL

	data, page, err := c.repo.GetListClassesPage(req, teacherId)
	if err != nil {
		return nil, err
	}

	for _, c := range data {
		class := classes_model.ClassHDRListDTOL{
			Id:           c.Id,
			Name:         c.Name,
			TotalStudent: c.TotalStudent,
		}

		if c.TotalModule > 0 {
			class.Progress = (float64(c.CompletedModule) / float64(c.TotalModule)) * 100
		}

		list = append(list, class)
	}

	response := map[string]interface{}{
		"data":       list,
		"page":       page.Page,
		"perpage":    page.Perpage,
		"total_page": page.TotalPage,
		"total_data": page.TotalData,
	}

	return &response, nil
}

func (c *classUsecase) CreateClass(req classes_model.ClassHdrDTO) error {
	model := classes_model.ClassHdr{
		Name:      req.Name,
		TeacherID: req.TeacherID,
		Status:    1,
		Level:     req.Level,
		Major:     req.Major,
	}

	if _, err := c.repo.InsertClassHDR(model); err != nil {
		return err
	}

	return nil
}

func (c *classUsecase) UpdateClass(req classes_model.ClassHdrDTO, StudentilePath string, CoursesFilePath string) error {
	model := classes_model.ClassHdr{
		ID:        req.ID,
		Name:      req.Name,
		TeacherID: req.TeacherID,
		Status:    req.Status,
		Level:     req.Level,
		Major:     req.Major,
	}

	if err := c.repo.UpdateClassHDR(model); err != nil {
		return err
	}

	if StudentilePath != "" {
		job := job_excel.ExcelJob{
			File:        &StudentilePath,
			ClassRepos:  c.repo,
			Importrepos: c.irepo,
			DB:          c.db,
			ClassId:     &req.ID,
			Type:        "student",
		}

		fmt.Println("enqueue excel job", StudentilePath)

		jobs.JobQueue <- job
	}

	if CoursesFilePath != "" {
		job := job_excel.ExcelJob{
			File:        &CoursesFilePath,
			ClassRepos:  c.repo,
			Importrepos: c.irepo,
			DB:          c.db,
			ClassId:     &req.ID,
			Type:        "course",
		}

		fmt.Println("enqueue excel job", CoursesFilePath)

		jobs.JobQueue <- job
	}

	return nil
}

func (c *classUsecase) GetAvailableStudent() ([]classes_model.StudentAvailable, error) {
	data, err := c.repo.GetAvailableStudents()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *classUsecase) GetClassByID(id uint64) (*classes_model.ClassHdr, error) {
	data, err := c.repo.GetClassById(classes_model.ClassHdr{ID: id})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *classUsecase) GetStudentClassByIDClass(id uint64) ([]classes_model.ClassStudent, error) {
	data, err := c.repo.GetStudentClass(classes_model.ClassStudent{ClassID: id})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *classUsecase) DeleteStudentClassByIDClass(id uint64) error {
	err := c.repo.DeleteStudentByIdCLass(id)
	if err != nil {
		return err
	}

	return nil
}

func (c *classUsecase) DeleteStudentClassByID(id uint64) error {
	err := c.repo.DeleteStudentClassById(classes_model.ClassStudent{ID: id})
	if err != nil {
		return err
	}

	return nil
}

func (c *classUsecase) GetAvailableCourseCLass() (*courses_model.AvailableCourse, error) {
	data, err := c.crepo.GetAvailableCourseClass()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *classUsecase) GetAllCourseClass(id uint64) ([]map[string]interface{}, error) {
	model := classes_model.ClassCourse{
		ClassID: id,
	}
	data, err := c.repo.GetAllClassCourses(model)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *classUsecase) GetInformDashboardClass(id int, teacherId int) (map[string]interface{}, error) {

	if id == 0 {
		return nil, fmt.Errorf("%w : %s", UnprocessableEntity, "Id kelas tidak boleh 0 atau kosong")
	}
	data, err := c.repo.GetInformDashboardClass(id, teacherId)

	if err != nil {
		return nil, fmt.Errorf("%w : %s", InternalServerError, err.Error())
	}

	modulAktifAny := data["modul_aktif"]
	modulTotalAny := data["modul_total"]

	modulAktif, ok1 := utils.ToFloat64(modulAktifAny)
	modulTotal, ok2 := utils.ToFloat64(modulTotalAny)

	var averageTotal float64

	if !ok1 {
		return nil, fmt.Errorf("%w : %s", InternalServerError, data["modul_aktif"])
	}

	if !ok2 {
		return nil, fmt.Errorf("%w : %s", InternalServerError, data["modul_aktif"])
	}

	if modulTotal != 0 {
		averageTotal = (float64(modulAktif) / float64(modulTotal)) * 100
	} else {
		averageTotal = 0
	}

	result := map[string]interface{}{
		"title":       data["name"],
		"courses":     data["courses"],
		"students":    data["students"],
		"assignments": data["assignments"],
		"average":     averageTotal,
	}

	return result, nil
}

func (c *classUsecase) GetAllModulByClassAndRole(classId int, teacherId int) ([]map[string]interface{}, error) {
	if classId == 0 {
		return nil, fmt.Errorf("%w : %s", UnprocessableEntity, "Class Id tidak boleh kosong")
	}

	data, err := c.repo.GetAllModulByClassAndRole(classId, teacherId)
	if err != nil {
		return nil, fmt.Errorf("%w : %s", InternalServerError, err.Error())
	}

	return data, nil
}

func (c *classUsecase) GetAvailableModulDash(classId int, teacherId int) ([]map[string]interface{}, error) {
	if classId == 0 {
		return nil, fmt.Errorf("%w : %s", UnprocessableEntity, "Class Id tidak boleh kosong")
	}

	data, err := c.repo.GetAvailableModulDash(classId, teacherId)
	if err != nil {
		return nil, fmt.Errorf("%w : %s", InternalServerError, err.Error())
	}

	return data, nil
}
