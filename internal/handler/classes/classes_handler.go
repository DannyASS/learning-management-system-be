package classes_handler

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	classes_model "github.com/DannyAss/users/internal/models/database_model/classes"
	classes_usecase "github.com/DannyAss/users/internal/usecase/classes"
	import_helper "github.com/DannyAss/users/pkg/import"
	"github.com/DannyAss/users/pkg/presentation"
	"github.com/DannyAss/users/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ClassesHandler struct {
	uc classes_usecase.IClassesUsecase
}

func NewClasssesHandler(uc classes_usecase.IClassesUsecase) *ClassesHandler {
	return &ClassesHandler{uc: uc}
}

func (ctrl *ClassesHandler) GetlistClassPage(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var params classes_model.ClassHDRRequest
		if ok := c.QueryParser(&params); ok != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(ok.Error()).
				Json(c)
		}

		userId := c.Locals("userID")
		teacherId, cOk := userId.(uint)
		if !cOk {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail("user ditolak").
				Json(c)
		}

		data, err := ctrl.uc.GetListClassesPage(params, int(teacherId))
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(data).
			Json(c)
	})
}

func (ctrl *ClassesHandler) CreateClassHdr(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var body classes_model.ClassHdrDTO
		if ok := c.BodyParser(&body); ok != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(ok.Error()).
				Json(c)
		}

		err := ctrl.uc.CreateClass(body)
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("sukses menyimpan data").
			Json(c)
	})
}

func (ctrl *ClassesHandler) UpdateClassHdr(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		id := c.Params("id")
		idInt, err1 := strconv.ParseUint(id, 10, 64)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}

		Studentfile, errFile := c.FormFile("student_file")

		if errFile != nil {
			Studentfile = nil
		}

		coursefile, errFile1 := c.FormFile("course_file")

		if errFile1 != nil {
			coursefile = nil
		}

		jsonData := c.FormValue("data")

		var body classes_model.ClassHdrDTO
		if ok := json.Unmarshal([]byte(jsonData), &body); ok != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(ok.Error()).
				Json(c)
		}

		if err := os.MkdirAll("./tmp", os.ModePerm); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed create tmp folder",
			})
		}

		var studentFilePath string
		var coursesFilePath string

		if Studentfile != nil {
			studentFilePath = "./tmp/" + Studentfile.Filename
			if err := c.SaveFile(Studentfile, studentFilePath); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "failed to save file"})
			}
		}

		if coursefile != nil {
			coursesFilePath = "./tmp/" + coursefile.Filename
			if err := c.SaveFile(coursefile, coursesFilePath); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "failed to save file"})
			}
		}

		body.ID = idInt
		err := ctrl.uc.UpdateClass(body, studentFilePath, coursesFilePath)
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("sukses memperbarui data").
			Json(c)
	})
}

func (ctrl *ClassesHandler) DownloadTemplate(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {

		data, err := ctrl.uc.GetAvailableStudent()
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error() + " " + "on populate student").
				Json(c)
		}

		buf, err2 := import_helper.GenerateStudentClassTemplate(data)
		if err2 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err2.Error() + " " + "on Generate template").
				Json(c)
		}

		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set("Content-Disposition", "attachment; filename=template_import_student.xlsx")

		return c.Send(buf.Bytes())
	})
}

func (ctrl *ClassesHandler) GetClassByIDClassHdr(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		id := c.Params("id")
		idUint, err1 := strconv.ParseUint(id, 10, 64)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}
		data, err := ctrl.uc.GetClassByID(idUint)
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(data).
			Json(c)
	})
}

func (ctrl *ClassesHandler) GetStudentClassByIDClass(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		id := c.Params("id")
		idUint, err1 := strconv.ParseUint(id, 10, 64)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}
		data, err := ctrl.uc.GetStudentClassByIDClass(idUint)
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(data).
			Json(c)
	})
}

func (ctrl *ClassesHandler) DeleteStudentClassByIDClass(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		id := c.Params("id")
		idUint, err1 := strconv.ParseUint(id, 10, 64)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}
		err := ctrl.uc.DeleteStudentClassByIDClass(idUint)
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Sukses menghapus data.").
			Json(c)
	})
}

func (ctrl *ClassesHandler) DeleteStudentClassByID(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		id := c.Params("id")
		idUint, err1 := strconv.ParseUint(id, 10, 64)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}
		err := ctrl.uc.DeleteStudentClassByID(idUint)
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Sukses menghapus data.").
			Json(c)
	})
}

func (ctrl *ClassesHandler) DownloadTemplateCourse(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {

		data, err := ctrl.uc.GetAvailableCourseCLass()
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error() + " " + "on populate student").
				Json(c)
		}

		buf, err2 := import_helper.GenerateCoursesTemplate(*data)
		if err2 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err2.Error() + " " + "on Generate template").
				Json(c)
		}

		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set("Content-Disposition", "attachment; filename=template_import_student.xlsx")

		return c.Send(buf.Bytes())
	})
}

func (ctrl *ClassesHandler) GetCOurseClassByIDClass(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		id := c.Params("id")
		idUint, err1 := strconv.ParseUint(id, 10, 64)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}
		data, err := ctrl.uc.GetAllCourseClass(idUint)
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(data).
			Json(c)
	})
}

func (ctrl *ClassesHandler) GetInformDasboardClass(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		id := c.Params("id")
		idUint, err1 := strconv.ParseUint(id, 10, 64)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}

		roleLocal := c.Locals("roleIDs")
		roleIds, ok := roleLocal.(uint)
		if !ok {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail("role ditolak").
				Json(c)
		}

		teacherId := 0

		if roleIds != 3 {
			userId, cek := c.Locals("userID").(uint)
			if !cek {
				return presentation.Response[any]().
					SetErrorCode("422").SetStatusCode(422).
					SetErrorDetail("user Id ditolak").
					Json(c)
			}

			teacherId = int(userId)
		}

		data, err := ctrl.uc.GetInformDashboardClass(int(idUint), teacherId)
		if err != nil {

			if errors.Is(classes_usecase.InternalServerError, err) {
				return presentation.Response[any]().
					SetErrorCode("500").SetStatusCode(500).
					SetErrorDetail(err.Error()).
					Json(c)
			} else {
				return presentation.Response[any]().
					SetErrorCode("422").SetStatusCode(422).
					SetErrorDetail(err.Error()).
					Json(c)
			}

		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(data).
			Json(c)
	})
}
