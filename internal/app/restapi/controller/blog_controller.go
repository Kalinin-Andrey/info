package controller

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"info/internal/domain/currency"
)

type blogController struct {
	router  *routing.Router
	service *currency.Service
}

func NewBlogController(router *routing.Router, service *currency.Service) *blogController {
	return &blogController{
		router:  router,
		service: service,
	}
}
