package presentation

import "github.com/gofiber/fiber/v2"

type (
	GlobalResponse[T any] struct {
		Status     bool   `json:"status"`
		StatusCode int    `json:"status_code"`
		StatusText string `json:"status_text,omitempty"`
		Message    string `json:"message,omitempty"`
		ErrDetail  string `json:"error_detail,omitempty"`
		ErrCode    string `json:"error_code,omitempty"`
		Data       T      `json:"data,omitempty"`
	}
)

func Response[T any]() *GlobalResponse[T] {
	return &GlobalResponse[T]{}
}

func (r *GlobalResponse[T]) SetStatus(status bool) *GlobalResponse[T] {
	r.Status = status
	return r
}

func (r *GlobalResponse[T]) SetMessage(msg string) *GlobalResponse[T] {
	r.Message = msg
	return r
}

func (r *GlobalResponse[T]) SetStatusCode(statusCode int) *GlobalResponse[T] {
	r.StatusCode = statusCode
	return r
}

func (r *GlobalResponse[T]) SetStatusText(statusText string) *GlobalResponse[T] {
	r.StatusText = statusText
	return r
}

func (r *GlobalResponse[T]) SetErrorDetail(detail string) *GlobalResponse[T] {
	r.ErrDetail = detail
	return r
}

func (r *GlobalResponse[T]) SetErrorCode(code string) *GlobalResponse[T] {
	r.ErrCode = code
	return r
}

func (r *GlobalResponse[T]) SetData(data T) *GlobalResponse[T] {
	r.Data = data
	return r
}

func (r *GlobalResponse[T]) Json(ctx *fiber.Ctx) error {
	return ctx.Status(r.StatusCode).JSON(r)
}
