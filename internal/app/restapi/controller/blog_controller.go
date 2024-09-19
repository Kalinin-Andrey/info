package controller

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"info/internal/domain/blog"
)

type blogController struct {
	router  *routing.Router
	service *blog.Service
}

func NewBlogController(router *routing.Router, service *blog.Service) *blogController {
	return &blogController{
		router:  router,
		service: service,
	}
}
