package job_excel

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/DannyAss/users/internal/database"
	classes_model "github.com/DannyAss/users/internal/models/database_model/classes"
	import_model "github.com/DannyAss/users/internal/models/database_model/import"
	classes_repository "github.com/DannyAss/users/internal/repositories/classes"
	courses_repository "github.com/DannyAss/users/internal/repositories/courses"
	import_repository "github.com/DannyAss/users/internal/repositories/import"
	"github.com/xuri/excelize/v2"
)

type ExcelJob struct {
	File        *string
	Importrepos import_repository.IImportRepository
	ClassRepos  classes_repository.IClassesRepository
	CourseRepos courses_repository.ICourseRepos
	DB          *database.DBManager
	ClassId     *uint64
	Type        string
}

func (j ExcelJob) Process() error {
	switch strings.ToLower(j.Type) {
	case "student":
		ImportStudent(j)
	case "course":
		ImportCourse(j)
	default:
		return errors.New("type job is not support")
	}

	return nil
}

func convertExcelToModel(path string, ClassId *uint64, repo classes_repository.IClassesRepository) ([]classes_model.ClassStudent, []error) {
	var Errors []error
	var ids []uint64
	f, err := excelize.OpenFile(path)
	if err != nil {
		Errors = append(Errors, err)
		return nil, Errors
	}
	defer f.Close()

	rows, err := f.GetRows("Template")
	if err != nil {
		Errors = append(Errors, err)
		return nil, Errors
	}

	// skip row pertama (header)
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		userId, err3 := strconv.ParseUint(safeGet(row, 0), 10, 64)
		if err3 != nil {
			Errors = append(Errors, err3)
		}
		ids = append(ids, userId)
	}

	checkIdUndifine, err5 := repo.GetExistingUserIDs(ids)

	if err5 != nil {
		Errors = append(Errors, err5)
		return nil, Errors
	}

	findingId := FindMissingIDs(ids, checkIdUndifine)
	if len(findingId) > 0 {
		for _, c := range findingId {
			er := errors.New("user id " + strconv.FormatUint(c, 10) + " tidak ditemukan")
			Errors = append(Errors, er)
			return nil, Errors
		}
	}

	getStudents, err4 := repo.GetNameStudents(ids, ClassId)
	if err4 != nil {
		Errors = append(Errors, err4)
	}

	if len(Errors) > 0 {
		return nil, Errors
	}

	return getStudents, nil
}

func safeGet(row []string, idx int) string {
	if idx < len(row) {
		return row[idx]
	}
	return ""
}

func ValidateStudent(data []classes_model.ClassStudent) []error {
	var errors []error

	for _, c := range data {
		required(c.ClassID, "uint64", &errors, "Class Id")
		required(c.UserID, "uint64", &errors, "User Id")
		required(c.Gender, "string", &errors, "Gander")
	}

	return errors
}

func required(data any, typedata string, err *[]error, field string) {
	switch strings.ToLower(typedata) {
	case "uint":
		v, ok := data.(uint64)
		if !ok || v == 0 {
			*err = append(*err, errors.New(field+" wajib diisi!"))
		}
	case "string":
		v, ok := data.(string)
		if !ok || strings.TrimSpace(v) == "" {
			*err = append(*err, errors.New(field+" wajib diisi!"))
		}
	default:
		*err = append(*err, errors.New("type data tidak valid"))
	}
}

func ConvertStudentIdToString(data []classes_model.ClassStudent) []string {
	var id []string
	for _, c := range data {
		idS := c.ID
		id = append(id, strconv.FormatUint(idS, 10))
	}
	return id
}

func ConvertErrorExist(data []classes_model.ClassStudent) []error {
	var errs []error
	for _, c := range data {
		idS := errors.New("siswa dengan user id " + strconv.FormatUint(c.UserID, 10) + "sudah terdaftar di kelas lain")
		errs = append(errs, idS)
	}
	return errs
}

func StartImport(file_name string, repos import_repository.IImportRepository, classId *uint64, mode string) (*int64, error) {
	var hdr import_model.ImportHdr
	switch strings.ToLower(mode) {
	case "course":
		hdr = import_model.ImportHdr{
			FileName: file_name,
			Type:     "courses class",
			Status:   "progress",
			ClassID:  classId,
		}

	case "student":
		hdr = import_model.ImportHdr{
			FileName: file_name,
			Type:     "student class",
			Status:   "progress",
			ClassID:  classId,
		}

	default:
		hdr = import_model.ImportHdr{
			FileName: file_name,
			Type:     "student class",
			Status:   "progress",
			ClassID:  classId,
		}

	}

	data, err := repos.InsertImportHDR(hdr)
	if err != nil {
		fmt.Println("error :", err)
		return nil, err
	}

	return &data.ID, nil
}

func UpdateImport(status string, onstep string, idImport int64, repo import_repository.IImportRepository, errors []error) error {
	updateImport := import_model.ImportHdr{
		ID:     idImport,
		Status: status,
		OnStep: &onstep,
	}

	err := repo.UpdateImportHDR(updateImport)
	if err != nil {
		return err
	}

	if status == "failed" {
		convertError := ConvertErrorToDetailImport(errors, idImport)
		err2 := repo.InsertImportDTL(convertError)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

func ConvertErrorToDetailImport(errs []error, import_id int64) []import_model.ImportDtl {
	var detailImport []import_model.ImportDtl
	for _, c := range errs {
		er := c.Error()
		dto := import_model.ImportDtl{
			ImportID: import_id,
			Desc:     &er,
		}

		detailImport = append(detailImport, dto)
	}

	return detailImport
}

func FindMissingIDs(input []uint64, exist []uint64) []uint64 {
	existMap := make(map[uint64]bool)

	for _, id := range exist {
		existMap[id] = true
	}

	var missing []uint64
	for _, id := range input {
		if !existMap[id] {
			missing = append(missing, id)
		}
	}

	return missing
}

func ImportStudent(j ExcelJob) error {
	fmt.Println("start job job")
	if j.DB == nil {
		panic("DBManager is nil")
	}

	if j.DB.GetDB() == nil {
		panic("gorm DB is nil")
	}
	var errs []error
	tx := j.DB.GetDB().Begin()

	OnStep := ""

	idImport, errImp := StartImport(*j.File, j.Importrepos, j.ClassId, "student")
	if errImp != nil {
		return errImp
	}

	repo := j.ClassRepos.WithTx(tx)
	defer func() {
		if len(errs) > 0 {
			fmt.Println("Error in step :", OnStep)
			tx.Rollback()
			UpdateImport("failed", OnStep, *idImport, j.Importrepos, errs)
		} else {
			tx.Commit()
			UpdateImport("success", "", *idImport, j.Importrepos, nil)
		}

	}()

	data, err := convertExcelToModel(*j.File, j.ClassId, repo)
	if err != nil {
		OnStep = "convert Excel"
		errs = append(errs, err...)
		return nil
	}

	ids := ConvertStudentIdToString(data)
	exist, err2 := repo.CekStudentIsExist(ids)
	if err2 != nil {
		OnStep = "Student Exist"
		existError := ConvertErrorExist(exist)
		errs = append(errs, existError...)
		return nil
	}

	err1 := repo.BulkInsertStudentClass(data)
	if err1 != nil {
		OnStep = "Insert student"
		errs = append(errs, err1)
		return nil
	}

	fmt.Println("start job job")

	return nil
}

func ImportCourse(j ExcelJob) error {
	fmt.Println("start job job")
	if j.DB == nil {
		panic("DBManager is nil")
	}

	if j.DB.GetDB() == nil {
		panic("gorm DB is nil")
	}
	var errs []error
	tx := j.DB.GetDB().Begin()

	OnStep := ""

	idImport, errImp := StartImport(*j.File, j.Importrepos, j.ClassId, "course")

	if errImp != nil {
		return errImp
	}

	repo := j.ClassRepos.WithTx(tx)

	defer func() {
		if len(errs) > 0 {
			fmt.Println("Error in step :", OnStep)
			tx.Rollback()
			UpdateImport("failed", OnStep, *idImport, j.Importrepos, errs)
		} else {
			tx.Commit()
			UpdateImport("success", "", *idImport, j.Importrepos, nil)
		}

	}()

	dataCourses, errs1 := convertExcelToModelCourse(*j.File, j.ClassId)

	if len(errs1) > 0 {
		OnStep = "convert excel to data courser"
		return nil
	}

	errCreate := repo.InsertBulkCourseClass(dataCourses)

	if errCreate != nil {
		errs = append(errs, errCreate)
		OnStep = "insert bulk course class"
		return nil
	}

	return nil
}

func ValidateCourses(data []classes_model.ClassCourse) []error {
	var errors []error

	for _, c := range data {
		required(c.CourseID, "uint64", &errors, "Courses Id")
		required(c.TeacherID, "uint64", &errors, "Teacher Id")
	}

	return errors
}

func convertExcelToModelCourse(path string, ClassId *uint64) ([]classes_model.ClassCourse, []error) {
	var Errors []error
	var ids []uint64
	var classCourses []classes_model.ClassCourse
	f, err := excelize.OpenFile(path)
	if err != nil {
		Errors = append(Errors, err)
		return nil, Errors
	}
	defer f.Close()

	rows, err := f.GetRows("Template")
	if err != nil {
		Errors = append(Errors, err)
		return nil, Errors
	}

	// skip row pertama (header)
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		dto := classes_model.ClassCourse{
			ClassID: *ClassId,
		}

		coursesId, err3 := strconv.ParseUint(safeGet(row, 0), 10, 64)
		userId, err4 := strconv.ParseUint(safeGet(row, 1), 10, 64)

		if err3 != nil {
			Errors = append(Errors, err3)
		}
		if err4 != nil {
			Errors = append(Errors, err4)
		}
		ids = append(ids, userId)

		dto.TeacherID = userId
		dto.CourseID = coursesId

		classCourses = append(classCourses, dto)
	}

	if len(Errors) > 0 {
		return nil, Errors
	}

	return classCourses, nil
}
