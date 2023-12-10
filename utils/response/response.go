package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code    int         `json:"code"`
	CodeMsg string      `json:"code_msg,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 统一封装成功响应值
func SuccessResponse(w http.ResponseWriter, resp interface{}) {
	var body Body

	body.Code = 0
	body.CodeMsg = "Success"
	body.Message = "Success"
	body.Data = resp

	httpx.OkJson(w, body)
}
