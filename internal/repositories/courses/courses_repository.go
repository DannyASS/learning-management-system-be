package courses_repository

import (
	"math"
	"strings"

	"github.com/DannyAss/users/internal/database"
	courses_model "github.com/DannyAss/users/internal/models/database_model/courses"
	user_model "github.com/DannyAss/users/internal/models/database_model/user"
	"gorm.io/gorm"
)

type coursesRepos struct {
	db *database.DBManager
	tx *gorm.DB
}

type ICourseRepos interface {
	WithTx(tx *gorm.DB) ICourseRepos
	InsertCourse(model courses_model.Course) (*courses_model.Course, error)
	UpdateData(filter courses_model.FilterCourse, model courses_model.Course) error
	Getlist(filter courses_model.FilterCourse, mode string, page *courses_model.CourseRequest) ([]courses_model.Course, *courses_model.CourseRequest, error)
	GetCourse(filter courses_model.FilterCourse, mode string) (*courses_model.Course, error)
	DeleteCourse(filter courses_model.FilterCourse) error
	GetCategory() ([]courses_model.Category, error)
	GetCategoryById(id int) (*courses_model.Category, error)
	GetAvailableCourseClass() (*courses_model.AvailableCourse, error)
}

func NewCourseRepos(db *database.DBManager) ICourseRepos {
	return &coursesRepos{db: db}
}

func (c *coursesRepos) WithTx(tx *gorm.DB) ICourseRepos {
	return &coursesRepos{db: c.db, tx: tx}
}

func (c *coursesRepos) getDB() *gorm.DB {
	if c.tx != nil {
		return c.tx
	}

	return c.db.GetDB()
}

func (c *coursesRepos) InsertCourse(model courses_model.Course) (*courses_model.Course, error) {
	tx := c.getDB()

	if err := tx.Create(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *coursesRepos) UpdateData(filter courses_model.FilterCourse, model courses_model.Course) error {
	tx := c.getDB()

	if err := tx.Where(&filter).Updates(&model).Error; err != nil {
		return err
	}

	return nil
}

func (c *coursesRepos) Getlist(filter courses_model.FilterCourse, mode string, page *courses_model.CourseRequest) ([]courses_model.Course, *courses_model.CourseRequest, error) {
	tx := c.getDB().Debug().Model(&courses_model.Course{})

	var (
		courses   []courses_model.Course
		totalData int64
	)

	if page != nil {
		mode = "like"
	}

	switch mode {
	case "full":
		tx = tx.Where(filter)
	case "like":
		title := filter.Title
		filter.Title = ""

		if title != "" {
			tx = tx.Where(" title LIKE ?", "%"+title+"%").Where(&filter)
		}

		if err1 := tx.Count(&totalData).Error; err1 != nil {
			return nil, nil, err1
		}

		if err2 := tx.Order(page.SortBy + " " + page.SortType).
			Limit(page.Perpage).
			Offset((page.Page - 1) * page.Perpage).
			Find(&courses).
			Error; err2 != nil {
			return nil, nil, err2
		}

		page.TotalPage = int(math.Ceil(float64(totalData) / float64(page.Perpage)))
		page.TotalData = int(totalData)

		return courses, page, nil
	default:
		tx = tx.Where(&filter)
	}

	if err3 := tx.Find(&courses).Error; err3 != nil {
		return nil, nil, err3
	}

	return courses, nil, nil
}

func (c *coursesRepos) GetCourse(filter courses_model.FilterCourse, mode string) (*courses_model.Course, error) {
	tx := c.getDB()

	var course courses_model.Course

	switch strings.ToLower(mode) {
	case "full":
		tx = tx.Where(filter)
	default:
		tx = tx.Where(filter)
	}

	if err := tx.Find(&course).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

func (c *coursesRepos) DeleteCourse(filter courses_model.FilterCourse) error {
	tx := c.getDB()

	if err := tx.Where(&filter).Delete(&courses_model.Course{}).Error; err != nil {
		return err
	}

	return nil
}

func (c *coursesRepos) GetCategory() ([]courses_model.Category, error) {
	tx := c.getDB()

	data := []courses_model.Category{}

	if err := tx.Table("categories").Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (c *coursesRepos) GetCategoryById(id int) (*courses_model.Category, error) {
	tx := c.getDB()

	data := courses_model.Category{}

	if err := tx.Table("categories").Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *coursesRepos) GetAvailableCourseClass() (*courses_model.AvailableCourse, error) {
	var (
		teacher []user_model.User
		courses []courses_model.Course
	)

	tx := c.getDB()
	tx2 := c.getDB()

	tx = tx.Raw(`Select
		a.id,
		a.name
		from
		users as a
		join user_roles as b on b.user_id = a.id
		where b.role_id = 1
	`)

	if err := tx.Scan(&teacher).Error; err != nil {
		return nil, err
	}

	tx2 = tx2.Raw(`
		Select
		a.id,
		a.title,
		a.description
		from courses a
		where a.status = 'published'
	`)

	if err2 := tx2.Scan(&courses).Error; err2 != nil {
		return nil, err2
	}

	result := courses_model.AvailableCourse{
		Teacher: teacher,
		Courses: courses,
	}
	return &result, nil
}
