package app

import (
	"net/http"

	"myGin/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

type ResponsePage struct {
	CommonResponseBody
	Page Pager `json:"pager,omitempty"`
}

type CommonResponseBody struct {
	RequestId      string      `json:"request_id"`
	Code           int         `json:"code"`
	Message        string      `json:"message"`
	Details        []string    `json:"details,omitempty"`
	EnglishMessage string      `json:"english_message,omitempty"`
	Data           interface{} `json:"data"`
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	r.Ctx.JSON(err.StatusCode(), CommonResponseBody{
		RequestId: GetRequestId(r.Ctx),
		Code:      err.Code(),
		Message:   err.Msg(),
		Details:   err.Details(),
	})
}

func (r *Response) ToResponse(data interface{}) {
	r.Ctx.JSON(http.StatusOK, CommonResponseBody{
		RequestId: GetRequestId(r.Ctx),
		Code:      errcode.Success.Code(),
		Message:   errcode.Success.Msg(),
		Data:      data,
	})
}

func (r *Response) ToResponseList(data interface{}, totalRows int) {
	responseData := ResponsePage{
		Page: Pager{
			Page:      GetPage(r.Ctx),
			PageSize:  GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	}
	responseData.Code = errcode.Success.Code()
	responseData.Message = errcode.Success.Msg()
	responseData.Data = data
	r.Ctx.JSON(http.StatusOK, responseData)
}