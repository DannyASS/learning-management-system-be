package import_helper

import (
	"bytes"
	"strconv"

	classes_model "github.com/DannyAss/users/internal/models/database_model/classes"
	courses_model "github.com/DannyAss/users/internal/models/database_model/courses"
	"github.com/xuri/excelize/v2"
)

func GenerateStudentClassTemplate(students []classes_model.StudentAvailable) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	// === SHEET 1 === Template
	sheet1 := "Template"
	f.SetSheetName("Sheet1", sheet1)

	headers := []string{"UserID"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet1, cell, h)
	}

	// styling header
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#DDDDDD"}, Pattern: 1},
	})
	f.SetCellStyle(sheet1, "A1", "D1", style)

	// === SHEET 2 === Daftar Student Available
	sheet2 := "StudentList"
	_, err2 := f.NewSheet(sheet2)
	if err2 != nil {
		return nil, err2
	}

	// header
	listHeader := []string{"UserID", "Name", "Gender"}
	for i, h := range listHeader {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet2, cell, h)
	}

	// data student
	for i, s := range students {
		row := i + 2

		f.SetCellValue(sheet2, "A"+strconv.Itoa(row), s.UserID)
		f.SetCellValue(sheet2, "B"+strconv.Itoa(row), s.Name)
		f.SetCellValue(sheet2, "C"+strconv.Itoa(row), s.Gender)
	}

	// styling
	f.SetCellStyle(sheet2, "A1", "C1", style)

	// set Sheet1 sebagai default
	f.SetActiveSheet(0)

	// convert ke buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func GenerateCoursesTemplate(courseClass courses_model.AvailableCourse) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	teacher := courseClass.Teacher
	courses := courseClass.Courses

	// === SHEET 1 === Template
	sheet1 := "Template"
	f.SetSheetName("Sheet1", sheet1)

	headers := []string{"Course Id", "Teacher Id"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet1, cell, h)
	}

	// styling header
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#DDDDDD"}, Pattern: 1},
	})
	f.SetCellStyle(sheet1, "A1", "B1", style)

	// === SHEET 2 === Daftar teacher Available
	sheet2 := "Teacher List"
	_, err2 := f.NewSheet(sheet2)
	if err2 != nil {
		return nil, err2
	}

	// header
	listHeader := []string{"Teacher Id", "Name"}
	for i, h := range listHeader {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet2, cell, h)
	}

	// data Teacher
	for i, s := range teacher {
		row := i + 2

		f.SetCellValue(sheet2, "A"+strconv.Itoa(row), s.ID)
		f.SetCellValue(sheet2, "B"+strconv.Itoa(row), s.Name)
	}

	// styling
	f.SetCellStyle(sheet2, "A1", "B1", style)

	// === SHEET 3 === Daftar Courses Available
	sheet3 := "Courses List"
	_, err3 := f.NewSheet(sheet3)
	if err3 != nil {
		return nil, err2
	}

	// header
	listHeaderCourses := []string{"Courses Id", "Title", "Description"}
	for i, h := range listHeaderCourses {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet3, cell, h)
	}

	// data Teacher
	for i, s := range courses {
		row := i + 2

		f.SetCellValue(sheet3, "A"+strconv.Itoa(row), s.ID)
		f.SetCellValue(sheet3, "B"+strconv.Itoa(row), s.Title)
		f.SetCellValue(sheet3, "C"+strconv.Itoa(row), *s.Description)
	}

	// styling
	f.SetCellStyle(sheet3, "A1", "C1", style)

	// set Sheet1 sebagai default
	f.SetActiveSheet(0)

	// convert ke buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func GenerateAbsenTemplate(students []classes_model.ClassStudent) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	// === SHEET 1 === Template
	sheet1 := "Template"
	f.SetSheetName("Sheet1", sheet1)

	headers := []string{"Student Id", "Absen", "Keterangan", "Tanggal"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet1, cell, h)
	}

	// styling header
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#DDDDDD"}, Pattern: 1},
	})
	f.SetCellStyle(sheet1, "A1", "D1", style)

	// === SHEET 2 === Daftar teacher Available
	sheet2 := "Keterangan"
	_, err2 := f.NewSheet(sheet2)
	if err2 != nil {
		return nil, err2
	}

	// header
	listHeader := []string{"Absen", "Keterangan"}
	for i, h := range listHeader {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet2, cell, h)
	}

	absen := map[string]string{
		"S": "Sakit",
		"A": "Alpa",
		"I": "Ijin",
	}

	row := 2
	// data Teacher
	for s := range absen {

		f.SetCellValue(sheet2, "A"+strconv.Itoa(row), s)
		f.SetCellValue(sheet2, "B"+strconv.Itoa(row), absen[s])
		row++
	}

	// styling
	f.SetCellStyle(sheet2, "A1", "B1", style)

	// === SHEET 3 === Daftar Courses Available
	sheet3 := "Students"
	_, err3 := f.NewSheet(sheet3)
	if err3 != nil {
		return nil, err3
	}

	// header
	listHeaderCourses := []string{"Student Id", "Student Name"}
	for i, h := range listHeaderCourses {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet3, cell, h)
	}

	// data Teacher
	for i, s := range students {
		row := i + 2

		f.SetCellValue(sheet3, "A"+strconv.Itoa(row), s.UserID)
		f.SetCellValue(sheet3, "B"+strconv.Itoa(row), s.StudentName)
	}

	// styling
	f.SetCellStyle(sheet3, "A1", "C1", style)

	// set Sheet1 sebagai default
	f.SetActiveSheet(0)

	// convert ke buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf, nil
}
