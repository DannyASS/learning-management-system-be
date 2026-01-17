package user_handler

import (
	"strconv"

	user_model "github.com/DannyAss/users/internal/models/database_model/user"
	user_usecase "github.com/DannyAss/users/internal/usecase/user"
	"github.com/DannyAss/users/pkg/presentation"
	"github.com/DannyAss/users/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUsecase user_usecase.IuserUsecae
}

func NewUserhandler(userUS user_usecase.IuserUsecae) *UserHandler {
	return &UserHandler{userUsecase: userUS}
}

var validate = validator.New()

func (cntrl *UserHandler) Login(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var request user_model.LoginDTO
		if err := c.BodyParser(&request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		token, refreshtoken, Rerr := cntrl.userUsecase.Login(request)

		if Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("401").SetStatusCode(401).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		c.Cookie(refreshtoken)

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(token).
			Json(c)
	})
}

func (cntrl *UserHandler) Register(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		var request user_model.RegisterDTO
		if err := c.BodyParser(&request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		if err := validate.Struct(request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("VALIDATION_ERROR").
				SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		if Rerr := cntrl.userUsecase.Register(request); Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").SetStatusCode(422).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Register Successfuly").
			Json(c)
	})
}

func (cntrl *UserHandler) RefreshToken(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {
		refreshToken := c.Cookies("refresh_token")
		if refreshToken == "" {
			return presentation.Response[any]().
				SetErrorCode("401").SetStatusCode(401).
				SetErrorDetail("refresh token is empty").
				Json(c)
		}

		token, Rerr := cntrl.userUsecase.RefreshToken(refreshToken)

		if Rerr != nil {
			return presentation.Response[any]().
				SetErrorCode("401").SetStatusCode(401).
				SetErrorDetail(Rerr.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(map[string]string{
				"token": token.Token,
			}).
			Json(c)
	})
}

func (cntrl *UserHandler) GetUsers(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {

		var request user_model.Pagination

		// Parse query param
		if err := c.QueryParser(&request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").
				SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		// Call usecase
		response, useErr := cntrl.userUsecase.GetUsers(request)
		if useErr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").
				SetStatusCode(422).
				SetErrorDetail(useErr.Error()).
				Json(c)
		}

		// Success response
		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(response).
			Json(c)
	})
}

func (cntrl *UserHandler) UpdateUser(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {

		var (
			request user_model.RequestUserGetList
		)

		Id := c.Params("id")

		NewId, err := strconv.ParseInt(Id, 10, 0)
		if err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").
				SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		// Parse query param
		if err := c.BodyParser(&request); err != nil {
			return presentation.Response[any]().
				SetErrorCode("422").
				SetStatusCode(422).
				SetErrorDetail(err.Error()).
				Json(c)
		}

		request.Id = int(NewId)

		// Call usecase
		useErr := cntrl.userUsecase.Updateuser(request)
		if useErr != nil {
			return presentation.Response[any]().
				SetErrorCode("422").
				SetStatusCode(422).
				SetErrorDetail(useErr.Error()).
				Json(c)
		}

		// Success response
		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetMessage("Success Update users").
			Json(c)
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "refresh token missing",
		})
	}

	// revoke on DB
	cookie, err := h.userUsecase.Logout(refreshToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to revoke token",
		})
	}

	// remove cookie
	c.Cookie(cookie)

	return presentation.Response[any]().
		SetStatus(true).
		SetStatusCode(200).
		SetMessage("Success Logout").
		Json(c)
}

func (h *UserHandler) GetAllTeacher(c *fiber.Ctx) error {
	return utils.TryCatch(c, func() error {

		data, err := h.userUsecase.GetAllTeacher()
		if err != nil {
			return presentation.Response[any]().
				SetStatus(false).
				SetStatusCode(422).
				SetMessage(err.Error()).
				Json(c)
		}

		return presentation.Response[any]().
			SetStatus(true).
			SetStatusCode(200).
			SetData(data).
			Json(c)
	})
}
