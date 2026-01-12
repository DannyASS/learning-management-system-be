package absen_handler

import (
	"errors"
	"strconv"

	absen_model "github.com/DannyAss/users/internal/models/database_model/absen"
	absen_usecase "github.com/DannyAss/users/internal/usecase/absen"
	"github.com/DannyAss/users/pkg/presentation"
	"github.com/DannyAss/users/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type (
	usecase = absen_usecase.IAbsenUsecase
	dtoReq  = absen_model.DTOAbsenRequest
)

type AbsenHandler struct {
	uc usecase
}

var (
	InternalServiceError = absen_usecase.InternalServiceError
	UnprocessableEntity  = absen_usecase.UnprocessableEntity
)

func NewHandler(uc usecase) *AbsenHandler {
	return &AbsenHandler{uc: uc}
}

func (cntrl *AbsenHandler) AbsenToday(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var dto dtoReq

		if ok := c.BodyParser(&dto); ok != nil {
			return presentation.Response[any]().
				SetStatus(false).
				SetStatusCode(422).
				SetMessage(ok.Error()).
				Json(c)
		}

		idClass := c.Params("id")

		id, ok2 := strconv.Atoi(idClass)
		if ok2 != nil {
			return presentation.Response[any]().
				SetStatus(false).
				SetStatusCode(422).
				SetMessage(ok2.Error()).
				Json(c)
		}

		err := cntrl.uc.AbsenToday(id, dto)

		if err != nil {
			if errors.Is(err, InternalServiceError) {
				return presentation.Response[any]().
					SetStatus(false).
					SetStatusCode(500).
					SetMessage(err.Error()).
					Json(c)
			} else {
				return presentation.Response[any]().
					SetStatus(false).
					SetStatusCode(422).
					SetMessage(err.Error()).
					Json(c)
			}
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Absen telah ditambahkan").
			Json(c)

	})
}

func (cntrl *AbsenHandler) GetAbsenToday(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		dto := c.Query("course_id")

		courseId, ok := strconv.Atoi(dto)

		if ok != nil {
			return presentation.Response[any]().
				SetStatus(false).
				SetStatusCode(422).
				SetMessage(ok.Error()).
				Json(c)
		}

		idClass := c.Params("id")

		id, ok2 := strconv.Atoi(idClass)
		if ok2 != nil {
			return presentation.Response[any]().
				SetStatus(false).
				SetStatusCode(422).
				SetMessage(ok2.Error()).
				Json(c)
		}

		absen := cntrl.uc.GetAbsen(id, courseId)

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(map[string]any{
				"absen": absen,
			}).
			Json(c)

	})
}
