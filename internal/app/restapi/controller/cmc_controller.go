package controller

import (
	"errors"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"info/internal/domain/currency"
	"info/internal/pkg/apperror"
	"info/internal/pkg/fasthttp_tools"
	"info/internal/pkg/log_key"
)

const (
	defaultLimit4Report = 10
)

type cmcController struct {
	logger  *zap.Logger
	router  *routing.Router
	service *currency.Service
}

func NewCmcController(logger *zap.Logger, router *routing.Router, service *currency.Service) *cmcController {
	return &cmcController{
		logger:  logger,
		router:  router,
		service: service,
	}
}

func (c *cmcController) Report_BiggestFall(rctx *routing.Context) (err error) {
	const metricName = "transitTariffController.Report_BiggestFall"
	var res *fasthttp_tools.Response

	limit, err := fasthttp_tools.ParseQueryArgUint(rctx.RequestCtx, "limit")
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			errMsg := "Parse params error "
			c.logger.Error(errMsg, zap.String(log_key.Func, metricName), zap.Error(err))
			res = fasthttp_tools.NewResponse_ErrBadRequest(errMsg + err.Error())
			fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusBadRequest, *res)
			return nil
		}
		limit = defaultLimit4Report
	}

	report, err := c.service.Report_BiggestFall(limit)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			errMsg := "Data for Report_BiggestFall was not found"
			c.logger.Error(errMsg, zap.String(log_key.Func, metricName), zap.Error(err))
			res = fasthttp_tools.NewResponse_ErrNotFound(errMsg)
			fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
			return nil
		}
		errMsg := "Failed to get Report_BiggestFall"
		c.logger.Error(errMsg, zap.String(log_key.Func, metricName), zap.Error(err))
		res = fasthttp_tools.NewResponse_ErrInternal()
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusInternalServerError, *res)
		return nil
	}

	res = fasthttp_tools.NewResponse_Success(report)
	fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusOK, *res)
	return nil
}

func (c *cmcController) Report_LongestFall(rctx *routing.Context) (err error) {
	const metricName = "transitTariffController.Report_BiggestFall"
	var res *fasthttp_tools.Response

	limit, err := fasthttp_tools.ParseQueryArgUint(rctx.RequestCtx, "limit")
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			errMsg := "Parse params error "
			c.logger.Error(errMsg, zap.String(log_key.Func, metricName), zap.Error(err))
			res = fasthttp_tools.NewResponse_ErrBadRequest(errMsg + err.Error())
			fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusBadRequest, *res)
			return nil
		}
		limit = defaultLimit4Report
	}

	report, err := c.service.Report_LongestFall(limit)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			errMsg := "Data for Report_BiggestFall was not found"
			c.logger.Error(errMsg, zap.String(log_key.Func, metricName), zap.Error(err))
			res = fasthttp_tools.NewResponse_ErrNotFound(errMsg)
			fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
			return nil
		}
		errMsg := "Failed to get Report_BiggestFall"
		c.logger.Error(errMsg, zap.String(log_key.Func, metricName), zap.Error(err))
		res = fasthttp_tools.NewResponse_ErrInternal()
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusInternalServerError, *res)
		return nil
	}

	res = fasthttp_tools.NewResponse_Success(report)
	fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusOK, *res)
	return nil
}
