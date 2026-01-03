package import_hanlder

import (
	"strconv"

	import_usecase "github.com/DannyAss/users/internal/usecase/import"
	"github.com/DannyAss/users/pkg/presentation"
	"github.com/DannyAss/users/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ImportHandler struct {
	uc import_usecase.IImportUsecase
}

func NewImportHandler(uc import_usecase.IImportUsecase) *ImportHandler {
	return &ImportHandler{uc: uc}
}

func (ctrl *ImportHandler) GetImportData(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		typeImport := c.Query("type")
		class_id := c.Query("id")
		uId, err := strconv.ParseUint(class_id, 10, 64)

		if err != nil {
			return presentation.Response[any]().
				SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		data, err2 := ctrl.uc.GetImport(typeImport, uId)

		if err2 != nil {
			return presentation.Response[any]().
				SetStatusCode(422).
				SetErrorDetail(err2.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatusCode(200).
			SetData(data).
			Json(c)
	})
}

func (ctrl *ImportHandler) GetTatusImportByClassId(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		typeImport := c.Query("type")
		class_id := c.Query("id")
		uId, err := strconv.ParseUint(class_id, 10, 64)

		if err != nil {
			return presentation.Response[any]().
				SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		data, err2 := ctrl.uc.GetImport(typeImport, uId)

		if err2 != nil {
			return presentation.Response[any]().
				SetStatusCode(422).
				SetErrorDetail(err2.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatusCode(200).
			SetData(data).
			Json(c)
	})
}
