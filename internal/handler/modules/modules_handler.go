package modules_handler

import (
	"strconv"

	modules_model "github.com/DannyAss/users/internal/models/database_model/modules"
	modules_usecase "github.com/DannyAss/users/internal/usecase/modules"
	"github.com/DannyAss/users/pkg/presentation"
	"github.com/DannyAss/users/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ModuleHandler struct {
	usecase modules_usecase.IModuleUsecase
}

func NewModuleHandler(uc modules_usecase.IModuleUsecase) *ModuleHandler {
	return &ModuleHandler{usecase: uc}
}

// ============================================================
// CREATE MODULE
// ============================================================
func (cntrl *ModuleHandler) CreateModule(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var request modules_model.ModuleDTO
		if err := c.BodyParser(&request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		if Rerr := cntrl.usecase.CreateModule(request); Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Create Module Successfully").
			Json(c)
	})
}

// ============================================================
// UPDATE MODULE
// ============================================================
func (cntrl *ModuleHandler) UpdateModule(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var request modules_model.ModuleDTO
		var filter modules_model.FilterModule

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

		filter.ID = int(Id)

		if Rerr := cntrl.usecase.UpdateModule(filter, request); Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Update Module Successfully").
			Json(c)
	})
}

// ============================================================
// GET LIST MODULE
// ============================================================
func (cntrl *ModuleHandler) GetListModule(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var request modules_model.ModuleRequest
		if err := c.QueryParser(&request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		response, Rerr := cntrl.usecase.GetListModule(request)
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

// ============================================================
// GET MODULE DETAIL
// ============================================================
func (cntrl *ModuleHandler) GetModule(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var filter modules_model.FilterModule

		Id, err1 := strconv.ParseInt(c.Params("id"), 10, 0)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}

		filter.ID = int(Id)

		response, Rerr := cntrl.usecase.GetModule(filter)
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

// ============================================================
// GET Courses
// ============================================================
func (cntrl *ModuleHandler) GetCourse(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {

		response, Rerr := cntrl.usecase.GetCourse()
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

// ============================================================
// DELETE MODULE
// ============================================================
func (cntrl *ModuleHandler) DeleteModule(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var filter modules_model.FilterModule

		Id, err1 := strconv.ParseInt(c.Params("id"), 10, 0)
		if err1 != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err1.Error()).
				Json(c)
		}

		filter.ID = int(Id)

		if Rerr := cntrl.usecase.DeleteModule(filter); Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Delete Module Successfully").
			Json(c)
	})
}
