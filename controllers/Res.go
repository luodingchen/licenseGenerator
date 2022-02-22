package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	InputError  = "Input error"
	PostError   = "Post error"
	DeleteError = "Delete error"
	PutError    = "Put error"
	GetError    = "Get error"

	PostSuccess   = "Post Success"
	DeleteSuccess = "Delete Success"
	PutSuccess    = "Put Success"
	GetSuccess    = "Get Success"
)

type Res struct {
	Code    int //0成功 500报错
	Data    interface{}
	Message string
	//	LicenseStatus string
	//	ResultMessage string
	//	Succeed bool
}

type ID struct {
	ID uint
}

type List struct {
	List interface{}
	//EndRow int
	//FirstPage int
	//HasNextPage bool
	//HasPreviousPage bool
	//IsFirstPage bool
	//IsLastPage bool
	//LastPage int
	//NavigateFirstPage int
	//NavigateLastPage int
	//NavigatePages int
	//NavigatePagesNums []int
	//NextPage int
	PageNum  int
	PageSize int
	//Pages int
	//PrePage int
	//Size int
	//StartRow int
	Total int
}

func NewResultMsg(ctx *gin.Context) *ResultMsg {
	return &ResultMsg{c: ctx}
}

type ResultMsg struct {
	c *gin.Context
}

func (r *ResultMsg) Success(msg string, data interface{}) {
	res := Res{}
	res.Message = msg
	res.Data = data
	r.c.JSON(http.StatusOK, res)
}

func (r *ResultMsg) Error(msg string) {
	res := Res{}
	res.Code = 500
	res.Message = msg
	res.Data = nil
	r.c.JSON(http.StatusOK, res)
}
