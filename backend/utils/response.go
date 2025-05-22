package utils

import (
	"encoding/json"
	"net/http"
)

// Response 通用响应结构
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondWithJSON 以JSON格式返回成功响应
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response := Response{
		Status: "success",
		Data:   payload,
	}

	RespondWithJSONDirect(w, code, response)
}

// RespondWithError 以JSON格式返回错误响应
func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := Response{
		Status:  "error",
		Message: message,
	}

	RespondWithJSONDirect(w, code, response)
}

// RespondWithJSONDirect 直接将数据以JSON格式响应
func RespondWithJSONDirect(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Total  int64       `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Data   interface{} `json:"data"`
}

// RespondWithPagination 返回分页数据响应
func RespondWithPagination(w http.ResponseWriter, code int, total int64, limit, offset int, data interface{}) {
	pagination := PaginationResponse{
		Total:  total,
		Limit:  limit,
		Offset: offset,
		Data:   data,
	}

	RespondWithJSON(w, code, pagination)
}
