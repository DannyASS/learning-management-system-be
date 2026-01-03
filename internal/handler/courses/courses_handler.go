package courses_handler

import (
	"strconv"

	courses_model "github.com/DannyAss/users/internal/models/database_model/courses"
	courses_usecase "github.com/DannyAss/users/internal/usecase/courses"
	"github.com/DannyAss/users/pkg/presentation"
	"github.com/DannyAss/users/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type CourseHandler struct {
	usecase courses_usecase.ICourseUsecase
}

func NewCourseHandler(uc courses_usecase.ICourseUsecase) *CourseHandler {
	return &CourseHandler{usecase: uc}
}

func (cntrl *CourseHandler) CreateCourse(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var request courses_model.CourseDTO
		if err := c.BodyParser(&request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		if Rerr := cntrl.usecase.CreateCourse(request); Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Create Course Successfuly").
			Json(c)
	})
}
func (cntrl *CourseHandler) UpdateCourse(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var request courses_model.CourseDTO
		var filter courses_model.FilterCourse
		if err := c.BodyParser(&request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		Id, err1 := strconv.ParseInt(c.Params("id"), 10, 0)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}

		filter.Id = int(Id)

		if Rerr := cntrl.usecase.UpdateCourse(filter, request); Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Upadte Course Successfuly").
			Json(c)
	})
}

func (cntrl *CourseHandler) GetListCourse(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var request courses_model.CourseRequest
		if err := c.QueryParser(&request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		response, Rerr := cntrl.usecase.GetListCourse(request)

		if Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(response).
			Json(c)
	})
}

func (cntrl *CourseHandler) GetCourse(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var filter courses_model.FilterCourse
		Id, err1 := strconv.ParseInt(c.Params("id"), 10, 0)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}

		filter.Id = int(Id)

		response, Rerr := cntrl.usecase.GetCourse(filter)

		if Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(response).
			Json(c)
	})
}
