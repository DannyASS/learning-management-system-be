package classes_repository

import (
	"errors"
	"math"

	"github.com/DannyAss/users/internal/database"
	classes_model "github.com/DannyAss/users/internal/models/database_model/classes"
	"gorm.io/gorm"
)

type classesRepository struct {
	db *database.DBManager
	tx *gorm.DB
}

type IClassesRepository interface {
	WithTx(tx *gorm.DB) IClassesRepository
	GetListClassesPage(page classes_model.ClassHDRRequest, teacherId int) ([]classes_model.ClassHDRQuery, *classes_model.ClassHDRRequest, error)
	BulkInsertStudentClass(model []classes_model.ClassStudent) error
	InsertStudentClass(model classes_model.ClassStudent) (*classes_model.ClassStudent, error)
	InsertClassHDR(model classes_model.ClassHdr) (*classes_model.ClassHdr, error)
	CekStudentIsExist(IdStudent []string) ([]classes_model.ClassStudent, error)
	UpdateClassHDR(model classes_model.ClassHdr) error
	GetAvailableStudents() ([]classes_model.StudentAvailable, error)
	GetClassById(model classes_model.ClassHdr) (*classes_model.ClassHdr, error)
	GetNameStudents(ids []uint64, classId *uint64) ([]classes_model.ClassStudent, error)
	GetExistingUserIDs(ids []uint64) ([]uint64, error)
	GetStudentClass(model classes_model.ClassStudent) ([]classes_model.ClassStudent, error)
	DeleteStudentClassById(model classes_model.ClassStudent) error
	DeleteStudentByIdCLass(id uint64) error
	InsertBulkCourseClass(model []classes_model.ClassCourse) error
	GetAllClassCourses(model classes_model.ClassCourse) ([]map[string]interface{}, error)
	GetInformDashboardClass(idClass int, teacherId int) (map[string]interface{}, error)
	GetAllModulByClassAndRole(idClass int, teacherId int) ([]map[string]interface{}, error)
	GetAvailableModulDash(idClass int, teacherId int) ([]map[string]interface{}, error)
}

func NewClassesRepos(db *database.DBManager) IClassesRepository {
	return &classesRepository{db: db}
}

func (c *classesRepository) WithTx(tx *gorm.DB) IClassesRepository {
	return &classesRepository{db: c.db, tx: tx}
}

func (c *classesRepository) getDB() *gorm.DB {
	if c.tx != nil {
		return c.tx
	}

	return c.db.GetDB()
}

func (c *classesRepository) GetListClassesPage(page classes_model.ClassHDRRequest, teacherId int) ([]classes_model.ClassHDRQuery, *classes_model.ClassHDRRequest, error) {
	var (
		data      []classes_model.ClassHDRQuery
		totalData int64
	)

	tx := c.getDB()

	moduleAktifSub := c.getDB().Select(`
		class_id,
		class_course_id,
		count(*) modules
	`).Table("class_module").
		Where("status = ?", "completed").
		Group("class_id, class_course_id")

	moduleTotalSub := c.getDB().Select(`
		class_id,
		class_course_id,
		count(*) modules
	`).Table("class_module").
		Group("class_id, class_course_id")

	courseSub := c.getDB().
		Select(`
			b1.class_id,
			count(b1.id) courses,
			b2.modules modul_aktif,
			b3.modules modul_total
		`).
		Table("class_course b1").
		Joins("join (?) b2 on b1.class_id = b2.class_id and b1.id = b2.class_course_id", moduleAktifSub).
		Joins("join (?) b3 on b1.class_id = b2.class_id and b1.id = b3.class_course_id", moduleTotalSub).
		Group("class_id, b2.modules, b3.modules")

	studentSub := c.getDB().
		Select(`
			class_id,
			count(*) students
		`).
		Table("class_student").
		Group("class_id")

	if teacherId != 0 {
		courseSub = courseSub.Where("teacher_id = ?", teacherId)
	}

	tx = tx.Select(`
		distinct
		a.id as id,
		a.name as name,
		ifnull(c.students, 0) total_student,
		b.modul_aktif completed_module,
		b.modul_total total_module

	`).Table("class_hdr a").
		Joins("left Join (?) b on b.class_id = a.id", courseSub).
		Joins("left Join (?) c on c.class_id = a.id", studentSub)

	if err1 := tx.Count(&totalData).Error; err1 != nil {
		return nil, nil, err1
	}

	if err := tx.Limit(page.Perpage).Offset((page.Page - 1) * page.Perpage).Scan(&data).Error; err != nil {
		return nil, nil, err
	}

	page.TotalPage = int(math.Ceil(float64(totalData) / float64(page.Perpage)))

	return data, &page, nil
}

func (c *classesRepository) InsertClassHDR(model classes_model.ClassHdr) (*classes_model.ClassHdr, error) {
	tx := c.getDB()

	if err := tx.Create(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *classesRepository) UpdateClassHDR(model classes_model.ClassHdr) error {
	tx := c.getDB()

	if err := tx.Updates(&model).Where("id = ?", model.ID).Error; err != nil {
		return err
	}

	return nil
}

func (c *classesRepository) DeleteClassHDR(model classes_model.ClassHdr) error {
	tx := c.getDB()

	if err := tx.Delete(&classes_model.ClassHdr{}).Where("id = ?", model.ID).Error; err != nil {
		return err
	}

	return nil
}

func (c *classesRepository) InsertStudentClass(model classes_model.ClassStudent) (*classes_model.ClassStudent, error) {
	tx := c.getDB()

	if err := tx.Create(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *classesRepository) BulkInsertStudentClass(model []classes_model.ClassStudent) error {
	tx := c.getDB()

	if err := tx.CreateInBatches(&model, 100).Error; err != nil {
		return err
	}

	return nil
}

func (c *classesRepository) CekStudentIsExist(IdStudent []string) ([]classes_model.ClassStudent, error) {
	var exist []classes_model.ClassStudent
	tx := c.getDB().Model(&classes_model.ClassStudent{})

	if err := tx.Where("id in ?", IdStudent).
		Where("status = 'active'").
		Find(&exist).Error; err != nil {
		return exist, errors.New("find user exist")
	}

	return nil, nil
}

func (c *classesRepository) GetAvailableStudents() ([]classes_model.StudentAvailable, error) {
	var result []classes_model.StudentAvailable

	// contoh: ambil student yang belum masuk kelas apa pun
	err := c.getDB().Raw(`
        SELECT u.id as user_id, u.name, u.gender 
        FROM users u
		Join user_roles ur on ur.user_id = u.id
        LEFT JOIN class_student cs ON u.id = cs.user_id and cs.status = 'active'
        WHERE cs.user_id IS NULL
		and ur.role_id = 2
		and u.is_verified = 1
    `).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *classesRepository) GetClassById(model classes_model.ClassHdr) (*classes_model.ClassHdr, error) {
	tx := c.getDB()

	if err := tx.Where("id = ?", model.ID).
		Preload("Students").
		Preload("Modules").
		Preload("ImportHdrs", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("id desc").Limit(1)
		}).
		First(&model).
		Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *classesRepository) GetNameStudents(ids []uint64, classId *uint64) ([]classes_model.ClassStudent, error) {
	var result []classes_model.StudentAvailable
	var students []classes_model.ClassStudent

	err := c.getDB().Raw(`
        SELECT u.id as user_id, u.name, u.gender 
        FROM users u
		WHERE u.id in ?
    `, ids).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	for _, c := range result {
		dto := classes_model.ClassStudent{
			UserID:      c.UserID,
			StudentName: c.Name,
			Gender:      c.Gender,
			Status:      "active",
			ClassID:     *classId,
		}

		students = append(students, dto)
	}

	return students, nil
}

func (r *classesRepository) GetExistingUserIDs(ids []uint64) ([]uint64, error) {
	var result []uint64

	err := r.getDB().
		Table("users").
		Select("id").
		Where("id IN ?", ids).
		Scan(&result).Error

	return result, err
}

func (c *classesRepository) DeleteStudentByIdCLass(id uint64) error {
	tx := c.getDB()

	if err := tx.Where("class_id = ?", id).Delete(&classes_model.ClassStudent{}).Error; err != nil {
		return err
	}

	return nil
}

func (c *classesRepository) GetStudentClass(model classes_model.ClassStudent) ([]classes_model.ClassStudent, error) {
	var (
		data []classes_model.ClassStudent
	)
	tx := c.getDB()

	if err := tx.Where("class_id = ?", model.ClassID).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (c *classesRepository) DeleteStudentClassById(model classes_model.ClassStudent) error {

	tx := c.getDB()

	if err := tx.Where("id = ?", model.ID).Delete(&classes_model.ClassStudent{}).Error; err != nil {
		return err
	}

	return nil

}

func (c *classesRepository) InsertBulkCourseClass(model []classes_model.ClassCourse) error {
	tx := c.getDB()

	if err := tx.CreateInBatches(&model, 100).Error; err != nil {
		return err
	}

	return nil
}

func (c *classesRepository) GetAllClassCourses(model classes_model.ClassCourse) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	tx := c.getDB()
	tx = tx.Raw(`
		Select
		distinct
		b.name,
		c.title
		from class_course as a
		join users b on a.teacher_id = b.id
		join courses c on c.id = a.course_id 
		where a.class_id = ?
	`, model.ClassID)
	if err := tx.Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (c *classesRepository) GetInformDashboardClass(idClass int, teacherId int) (map[string]interface{}, error) {
	var data map[string]interface{}

	tx := c.getDB()

	moduleAktifSub := c.getDB().Select(`
		class_id,
		class_course_id,
		count(*) modules
	`).Table("class_module").
		Where("status = ?", "completed").
		Group("class_id, class_course_id")

	moduleTotalSub := c.getDB().Select(`
		class_id,
		class_course_id,
		count(*) modules
	`).Table("class_module").
		Group("class_id, class_course_id")

	courseSub := c.getDB().
		Select(`
			b1.class_id,
			count(b1.id) courses,
			b2.modules modul_aktif,
			b3.modules modul_total
		`).
		Table("class_course b1").
		Joins("join (?) b2 on b1.class_id = b2.class_id and b1.id = b2.class_course_id", moduleAktifSub).
		Joins("join (?) b3 on b1.class_id = b2.class_id and b1.id = b3.class_course_id", moduleTotalSub).
		Group("class_id, b2.modules, b3.modules")

	studentSub := c.getDB().
		Select(`
			class_id,
			count(*) students
		`).
		Table("class_student").
		Group("class_id")

	assigmentSub := c.getDB().
		Select(`
			class_id,
			count(*) assignments
		`).
		Model(&classes_model.ClassAssignment{}).
		Where("status = ?", "active").
		Group("class_id")

	if teacherId != 0 {
		courseSub = courseSub.Where("teacher_id = ?", teacherId)
	}

	tx = tx.Select(`
		ifnull(b.courses, 0) courses,
		ifnull(c.students, 0) students,
		ifnull(d.assignments, 0) assignments,
		ifnull(b.modul_aktif, 0) modul_aktif,
		ifnull(b.modul_total, 0) modul_total

	`).Table("class_hdr a").
		Joins("left Join (?) b on b.class_id = a.id", courseSub).
		Joins("left Join (?) c on c.class_id = a.id", studentSub).
		Joins("left Join (?) d on d.class_id = a.id", assigmentSub).Where("a.id = ?", idClass).Scan(&data)

	if err := tx.Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *classesRepository) GetAllModulByClassAndRole(idClass int, teacherId int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}

	tx := r.getDB().Table("class_course a")

	tx = tx.Select(`
		b.id,
		c.title as courses,
		d.title name,
		b.status
	`)

	tx = tx.Joins(`join class_module b on b.class_id = a.class_id and b.class_course_id = a.id`)
	tx = tx.Joins(`join courses c on c.id = a.course_id`)
	tx = tx.Joins(`join course_modules d on d.id = b.module_id`)

	tx = tx.Where("a.class_id = ?", idClass)

	if teacherId != 0 {
		tx = tx.Where("a.teacher_id = ?", teacherId)
	}

	if err := tx.Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *classesRepository) GetAvailableModulDash(idClass int, teacherId int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}

	tx := r.getDB().Table("class_course a")

	courseSUb := r.getDB().Select("db.module_id").Table("class_course da").
		Joins("join class_module db on da.class_id = db.class_id and db.class_course_id = da.id").
		Where("da.class_id = a.class_id").
		Where("da.course_id = a.course_id")

	tx = tx.Select(`
		b.title as course,
		c.title as module
	`).
		Joins("join courses b on a.course_id = b.id").
		Joins("join course_modules c on c.course_id = a.course_id").
		Where("c.id not in (?)", courseSUb).
		Where("a.teacher_id = ?", teacherId).
		Where("a.class_id = ?", idClass)

	if err := tx.Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
